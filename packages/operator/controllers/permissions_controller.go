package controllers

import (
	"context"
	"fmt"
	"strings"

	sharekubev1alpha1 "github.com/miloszsobczak/sharekube/packages/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PermissionsManager handles dynamic RBAC for ShareKube resources
type PermissionsManager struct {
	client client.Client
	scheme *runtime.Scheme
}

// NewPermissionsManager creates a new permissions manager
func NewPermissionsManager(client client.Client, scheme *runtime.Scheme) *PermissionsManager {
	return &PermissionsManager{
		client: client,
		scheme: scheme,
	}
}

// ResourcePermission holds information about Kubernetes resource types
type ResourcePermission struct {
	APIGroup  string
	Resources []string
}

// getResourceMapping returns the mapping of resource kinds to their API groups
func getResourceMapping() map[string]ResourcePermission {
	return map[string]ResourcePermission{
		"Deployment": {
			APIGroup:  "apps",
			Resources: []string{"deployments"},
		},
		"StatefulSet": {
			APIGroup:  "apps",
			Resources: []string{"statefulsets"},
		},
		"DaemonSet": {
			APIGroup:  "apps",
			Resources: []string{"daemonsets"},
		},
		"ReplicaSet": {
			APIGroup:  "apps",
			Resources: []string{"replicasets"},
		},
		"Service": {
			APIGroup:  "",
			Resources: []string{"services"},
		},
		"ConfigMap": {
			APIGroup:  "",
			Resources: []string{"configmaps"},
		},
		"Secret": {
			APIGroup:  "",
			Resources: []string{"secrets"},
		},
		"Pod": {
			APIGroup:  "",
			Resources: []string{"pods"},
		},
		"PersistentVolumeClaim": {
			APIGroup:  "",
			Resources: []string{"persistentvolumeclaims"},
		},
		"Job": {
			APIGroup:  "batch",
			Resources: []string{"jobs"},
		},
		"CronJob": {
			APIGroup:  "batch",
			Resources: []string{"cronjobs"},
		},
		"Ingress": {
			APIGroup:  "networking.k8s.io",
			Resources: []string{"ingresses"},
		},
		"NetworkPolicy": {
			APIGroup:  "networking.k8s.io",
			Resources: []string{"networkpolicies"},
		},
	}
}

// EnsurePermissions ensures that the necessary permissions are created for a ShareKube resource
func (pm *PermissionsManager) EnsurePermissions(ctx context.Context, sharekube *sharekubev1alpha1.ShareKube) error {
	logger := log.FromContext(ctx)
	logger.Info("Ensuring permissions for ShareKube resource", "namespace", sharekube.Namespace, "name", sharekube.Name)

	// Skip dynamic permission creation if not enabled or explicitly disabled
	if sharekube.Spec.AccessControl == nil || !sharekube.Spec.AccessControl.Restrict {
		logger.Info("Dynamic permissions not enabled for this ShareKube resource")
		return nil
	}

	// Map to track required permissions
	requiredPermissions := make(map[string][]string)

	// Collect resource types from ShareKube resources list
	for _, resource := range sharekube.Spec.Resources {
		resourceMap := getResourceMapping()
		resourceInfo, exists := resourceMap[resource.Kind]
		if !exists {
			logger.Info("Unknown resource kind", "kind", resource.Kind)
			continue
		}

		if _, ok := requiredPermissions[resourceInfo.APIGroup]; !ok {
			requiredPermissions[resourceInfo.APIGroup] = []string{}
		}

		for _, resourceType := range resourceInfo.Resources {
			if !contains(requiredPermissions[resourceInfo.APIGroup], resourceType) {
				requiredPermissions[resourceInfo.APIGroup] = append(requiredPermissions[resourceInfo.APIGroup], resourceType)
			}
		}
	}

	// Create or update source namespace roles
	sourceNamespace := sharekube.Namespace
	sourceRoleErr := pm.createOrUpdateRole(ctx, sharekube, sourceNamespace, requiredPermissions, false)
	if sourceRoleErr != nil {
		return fmt.Errorf("failed to create source roles: %w", sourceRoleErr)
	}

	// Create or update target namespace roles
	targetNamespace := sharekube.Spec.TargetNamespace
	targetRoleErr := pm.createOrUpdateRole(ctx, sharekube, targetNamespace, requiredPermissions, true)
	if targetRoleErr != nil {
		return fmt.Errorf("failed to create target roles: %w", targetRoleErr)
	}

	return nil
}

// createOrUpdateRole creates or updates a Role in the specified namespace
func (pm *PermissionsManager) createOrUpdateRole(
	ctx context.Context,
	sharekube *sharekubev1alpha1.ShareKube,
	namespace string,
	requiredPermissions map[string][]string,
	isTarget bool,
) error {
	logger := log.FromContext(ctx)

	// Check if namespace exists, create if it doesn't (for target namespace)
	if isTarget {
		var ns corev1.Namespace
		err := pm.client.Get(ctx, types.NamespacedName{Name: namespace}, &ns)
		if err != nil && errors.IsNotFound(err) {
			newNs := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			}
			if err := pm.client.Create(ctx, newNs); err != nil {
				return fmt.Errorf("failed to create target namespace: %w", err)
			}
			logger.Info("Created target namespace", "namespace", namespace)
		} else if err != nil {
			return fmt.Errorf("failed to check target namespace: %w", err)
		}
	}

	// Create policy rules from required permissions
	var rules []rbacv1.PolicyRule
	for apiGroup, resources := range requiredPermissions {
		rule := rbacv1.PolicyRule{
			APIGroups: []string{apiGroup},
			Resources: resources,
		}

		if isTarget {
			rule.Verbs = []string{"get", "list", "watch", "create", "update", "patch", "delete"}
		} else {
			rule.Verbs = []string{"get", "list", "watch"}
		}

		rules = append(rules, rule)
	}

	// Add permission for finalizers and status if this is the source namespace
	if !isTarget {
		rules = append(rules, rbacv1.PolicyRule{
			APIGroups: []string{"sharekube.dev"},
			Resources: []string{"sharekubes", "sharekubes/status", "sharekubes/finalizers"},
			Verbs:     []string{"get", "list", "watch", "update", "patch"},
		})
	}

	// Create or update Role
	roleSuffix := "source"
	if isTarget {
		roleSuffix = "target"
	}

	roleName := fmt.Sprintf("sharekube-%s-%s", sharekube.Name, roleSuffix)

	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, pm.client, role, func() error {
		role.Rules = rules

		// Set owner reference only if it's in the same namespace
		if namespace == sharekube.Namespace {
			if err := controllerutil.SetControllerReference(sharekube, role, pm.scheme); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create/update role: %w", err)
	}

	// Now create or update RoleBinding
	serviceAccount := "sharekube-controller-manager"
	serviceAccountNamespace := "sharekube-system"

	binding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName + "-binding",
			Namespace: namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(ctx, pm.client, binding, func() error {
		binding.RoleRef = rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     roleName,
		}

		binding.Subjects = []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccount,
				Namespace: serviceAccountNamespace,
			},
		}

		// Set owner reference only if it's in the same namespace
		if namespace == sharekube.Namespace {
			if err := controllerutil.SetControllerReference(sharekube, binding, pm.scheme); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create/update role binding: %w", err)
	}

	// Update ShareKube status to track the permissions
	permissionRef := fmt.Sprintf("%s/%s", namespace, roleName)
	if !contains(sharekube.Status.DynamicPermissions, permissionRef) {
		sharekube.Status.DynamicPermissions = append(sharekube.Status.DynamicPermissions, permissionRef)
	}

	logger.Info("Created/updated role and binding", "namespace", namespace, "role", roleName)

	return nil
}

// CleanupPermissions removes the dynamic permissions for a ShareKube resource
func (pm *PermissionsManager) CleanupPermissions(ctx context.Context, sharekube *sharekubev1alpha1.ShareKube) error {
	// For permissions in other namespaces that don't have owner references
	for _, permRef := range sharekube.Status.DynamicPermissions {
		parts := strings.Split(permRef, "/")
		if len(parts) != 2 {
			continue
		}

		namespace := parts[0]
		roleName := parts[1]

		// Skip if it's in the same namespace (will be cleaned up by garbage collection)
		if namespace == sharekube.Namespace {
			continue
		}

		// Delete role binding
		binding := &rbacv1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleName + "-binding",
				Namespace: namespace,
			},
		}

		if err := pm.client.Delete(ctx, binding); err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("failed to delete role binding: %w", err)
		}

		// Delete role
		role := &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      roleName,
				Namespace: namespace,
			},
		}

		if err := pm.client.Delete(ctx, role); err != nil && !errors.IsNotFound(err) {
			return fmt.Errorf("failed to delete role: %w", err)
		}
	}

	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}
