package controllers

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	sharekubev1alpha1 "github.com/miloszsobczak/sharekube/packages/operator/api/v1alpha1"
	"github.com/miloszsobczak/sharekube/packages/operator/pkg/resources"
)

// ShareKubeReconciler reconciles a ShareKube object
type ShareKubeReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Config    *rest.Config
	DynClient dynamic.Interface
}

//+kubebuilder:rbac:groups=sharekube.dev,resources=sharekubes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sharekube.dev,resources=sharekubes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sharekube.dev,resources=sharekubes/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create
//+kubebuilder:rbac:groups=core,resources=services;configmaps;secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// The ShareKubeFinalizer is used to clean up resources when a ShareKube resource is deleted
const ShareKubeFinalizer = "sharekube.dev/finalizer"

// Reconcile handles the main reconciliation loop for ShareKube resources
func (r *ShareKubeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling ShareKube", "Request.Namespace", req.Namespace, "Request.Name", req.Name)

	// Fetch the ShareKube instance
	sharekube := &sharekubev1alpha1.ShareKube{}
	err := r.Get(ctx, req.NamespacedName, sharekube)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Object not found, return
			logger.Info("ShareKube resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request
		logger.Error(err, "Failed to get ShareKube")
		return ctrl.Result{}, err
	}

	// Initialize status if it's a new resource
	if sharekube.Status.Phase == "" {
		now := metav1.Now()
		sharekube.Status.Phase = "Initializing"
		sharekube.Status.CreationTime = &now

		// Calculate expiration time based on TTL
		ttlDuration, err := time.ParseDuration(sharekube.Spec.TTL)
		if err != nil {
			logger.Error(err, "Invalid TTL format", "TTL", sharekube.Spec.TTL)
			sharekube.Status.Phase = "Error"
			if err := r.Status().Update(ctx, sharekube); err != nil {
				logger.Error(err, "Failed to update ShareKube status")
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		}

		expirationTime := metav1.NewTime(now.Add(ttlDuration))
		sharekube.Status.ExpirationTime = &expirationTime

		if err := r.Status().Update(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to update ShareKube status")
			return ctrl.Result{}, err
		}
	}

	// Add finalizer if not present
	if !controllerutil.ContainsFinalizer(sharekube, ShareKubeFinalizer) {
		controllerutil.AddFinalizer(sharekube, ShareKubeFinalizer)
		if err := r.Update(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Check if resource is being deleted
	if !sharekube.DeletionTimestamp.IsZero() {
		return r.handleDeletion(ctx, sharekube)
	}

	// Check if TTL has expired
	if sharekube.Status.ExpirationTime != nil && sharekube.Status.ExpirationTime.Before(&metav1.Time{Time: time.Now()}) {
		logger.Info("TTL expired, deleting ShareKube resource")
		if err := r.Delete(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to delete expired ShareKube resource")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Ensure target namespace exists
	targetNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: sharekube.Spec.TargetNamespace,
		},
	}

	if err := r.Get(ctx, types.NamespacedName{Name: targetNamespace.Name}, targetNamespace); err != nil {
		if apierrors.IsNotFound(err) {
			// Create the namespace if it doesn't exist
			logger.Info("Creating target namespace", "Namespace", targetNamespace.Name)
			if err := r.Create(ctx, targetNamespace); err != nil {
				logger.Error(err, "Failed to create target namespace")
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "Failed to get target namespace")
			return ctrl.Result{}, err
		}
	}

	// Update status to Processing if still Initializing
	if sharekube.Status.Phase == "Initializing" {
		sharekube.Status.Phase = "Processing"
		if err := r.Status().Update(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to update ShareKube status")
			return ctrl.Result{}, err
		}
	}

	// Process resources to copy
	copiedResources, err := r.processResources(ctx, sharekube)
	if err != nil {
		logger.Error(err, "Failed to process resources")
		sharekube.Status.Phase = "Error"
		if err := r.Status().Update(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to update ShareKube status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, err
	}

	// Update status with copied resources
	sharekube.Status.CopiedResources = copiedResources
	sharekube.Status.Phase = "Ready"
	if err := r.Status().Update(ctx, sharekube); err != nil {
		logger.Error(err, "Failed to update ShareKube status")
		return ctrl.Result{}, err
	}

	// Requeue to check TTL expiration
	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

// processResources copies the specified resources from source to target namespace
func (r *ShareKubeReconciler) processResources(ctx context.Context, sharekube *sharekubev1alpha1.ShareKube) ([]string, error) {
	logger := log.FromContext(ctx)
	var copiedResources []string

	// Create owner reference for all copied resources
	ownerRef := metav1.OwnerReference{
		APIVersion:         sharekube.APIVersion,
		Kind:               sharekube.Kind,
		Name:               sharekube.Name,
		UID:                sharekube.UID,
		Controller:         &[]bool{true}[0],
		BlockOwnerDeletion: &[]bool{true}[0],
	}

	// Create resource handler with owner reference and ShareKube info
	resourceHandler := resources.NewResourceHandler(
		r.Client,
		r.DynClient,
		r.Scheme,
		ownerRef,
		sharekube.Name,
		sharekube.Namespace,
	)

	for _, resource := range sharekube.Spec.Resources {
		resourceNamespace := resource.Namespace
		if resourceNamespace == "" {
			resourceNamespace = sharekube.Namespace
		}

		logger.Info("Copying resource",
			"Kind", resource.Kind,
			"Name", resource.Name,
			"SourceNamespace", resourceNamespace,
			"TargetNamespace", sharekube.Spec.TargetNamespace)

		// Use the resource handler to copy the resource
		err := resourceHandler.CopyResource(ctx, resource.Kind, resource.Name, resourceNamespace, sharekube.Spec.TargetNamespace)
		if err != nil {
			logger.Error(err, "Failed to copy resource",
				"Kind", resource.Kind,
				"Name", resource.Name,
				"SourceNamespace", resourceNamespace)
			continue
		}

		resourceRef := fmt.Sprintf("%s/%s/%s", resource.Kind, resourceNamespace, resource.Name)
		copiedResources = append(copiedResources, resourceRef)
	}

	return copiedResources, nil
}

// handleDeletion performs cleanup when a ShareKube resource is being deleted
func (r *ShareKubeReconciler) handleDeletion(ctx context.Context, sharekube *sharekubev1alpha1.ShareKube) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Handling deletion of ShareKube resource")

	// Check if finalizer is present
	if controllerutil.ContainsFinalizer(sharekube, ShareKubeFinalizer) {
		// We can't use Kubernetes owner references for cross-namespace cleanup
		// so we must explicitly clean up resources with our labels
		if err := r.cleanupResources(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to clean up resources")
			return ctrl.Result{}, err
		}

		// Remove finalizer to allow deletion
		controllerutil.RemoveFinalizer(sharekube, ShareKubeFinalizer)
		if err := r.Update(ctx, sharekube); err != nil {
			logger.Error(err, "Failed to remove finalizer")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// cleanupResources removes resources that were created by this ShareKube
func (r *ShareKubeReconciler) cleanupResources(ctx context.Context, sharekube *sharekubev1alpha1.ShareKube) error {
	logger := log.FromContext(ctx)
	logger.Info("Cleaning up resources", "TargetNamespace", sharekube.Spec.TargetNamespace)

	// Create Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(r.Config)
	if err != nil {
		logger.Error(err, "Failed to create Kubernetes clientset")
		return err
	}

	// Define the label selector to match resources created by this ShareKube instance
	// We don't need the 'copied' label since 'source-namespace' already implies it was copied
	labelSelector := fmt.Sprintf(
		"sharekube.dev/owner-name=%s,sharekube.dev/owner-namespace=%s",
		sharekube.Name,
		sharekube.Namespace,
	)

	// Delete deployments
	if err := clientset.AppsV1().Deployments(sharekube.Spec.TargetNamespace).DeleteCollection(
		ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labelSelector}); err != nil {
		logger.Error(err, "Failed to delete Deployments")
	}

	// Delete services
	services, err := clientset.CoreV1().Services(sharekube.Spec.TargetNamespace).List(
		ctx, metav1.ListOptions{LabelSelector: labelSelector})
	if err == nil {
		for _, svc := range services.Items {
			if err := clientset.CoreV1().Services(sharekube.Spec.TargetNamespace).Delete(
				ctx, svc.Name, metav1.DeleteOptions{}); err != nil {
				logger.Error(err, "Failed to delete Service", "Name", svc.Name)
			}
		}
	} else {
		logger.Error(err, "Failed to list Services")
	}

	// Delete configmaps
	if err := clientset.CoreV1().ConfigMaps(sharekube.Spec.TargetNamespace).DeleteCollection(
		ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labelSelector}); err != nil {
		logger.Error(err, "Failed to delete ConfigMaps")
	}

	// Delete secrets
	if err := clientset.CoreV1().Secrets(sharekube.Spec.TargetNamespace).DeleteCollection(
		ctx, metav1.DeleteOptions{}, metav1.ListOptions{LabelSelector: labelSelector}); err != nil {
		logger.Error(err, "Failed to delete Secrets")
	}

	// Use dynamic client to delete any other resources that we might have created
	// For brevity, we're omitting this, but in a real implementation you would use
	// discovery to find all installed types and then delete those with our labels

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShareKubeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sharekubev1alpha1.ShareKube{}).
		Complete(r)
}
