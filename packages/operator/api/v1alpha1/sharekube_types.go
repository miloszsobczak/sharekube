package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Resource defines a Kubernetes resource to copy
type Resource struct {
	// Kind is the type of Kubernetes resource (e.g., Deployment, Service)
	Kind string `json:"kind"`
	
	// Name is the name of the resource to copy
	Name string `json:"name"`
	
	// Namespace is the source namespace (optional, defaults to ShareKube CRD namespace)
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// TransformationRule defines how resources should be transformed during copying
// This is a future feature and not implemented in the initial version
type TransformationRule struct {
	// Kind is the resource type to apply transformations to
	Kind string `json:"kind"`
	
	// RemoveFields is a list of fields to remove from the resource
	// +optional
	RemoveFields []string `json:"removeFields,omitempty"`
}

// TargetCluster defines a remote Kubernetes cluster
// This is a future feature and not implemented in the initial version
type TargetCluster struct {
	// Name of the target cluster
	Name string `json:"name"`
	
	// KubeconfigSecret is the name of the secret containing the kubeconfig
	KubeconfigSecret string `json:"kubeconfigSecret"`
}

// ShareKubeSpec defines the desired state of ShareKube
type ShareKubeSpec struct {
	// TargetNamespace is the destination namespace for copied resources
	TargetNamespace string `json:"targetNamespace"`
	
	// TTL is the time-to-live for the preview environment (e.g., 1h, 24h, 7d)
	TTL string `json:"ttl"`
	
	// Resources is the list of resources to be copied
	Resources []Resource `json:"resources"`
	
	// TransformationRules is the list of transformation rules to apply (future feature)
	// +optional
	TransformationRules []TransformationRule `json:"transformationRules,omitempty"`
	
	// TargetCluster is the configuration for a remote cluster (future feature)
	// +optional
	TargetCluster *TargetCluster `json:"targetCluster,omitempty"`
}

// ShareKubeStatus defines the observed state of ShareKube
type ShareKubeStatus struct {
	// Phase is the current phase of the ShareKube resource
	// +optional
	Phase string `json:"phase,omitempty"`
	
	// CreationTime is when the preview environment was created
	// +optional
	CreationTime *metav1.Time `json:"creationTime,omitempty"`
	
	// ExpirationTime is when the preview environment will be deleted
	// +optional
	ExpirationTime *metav1.Time `json:"expirationTime,omitempty"`
	
	// CopiedResources is the list of resources that were successfully copied
	// +optional
	CopiedResources []string `json:"copiedResources,omitempty"`
	
	// Conditions represent the latest available observations of the ShareKube's state
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Target",type="string",JSONPath=".spec.targetNamespace"
//+kubebuilder:printcolumn:name="TTL",type="string",JSONPath=".spec.ttl"
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ShareKube is the Schema for the sharekubes API
type ShareKube struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ShareKubeSpec   `json:"spec,omitempty"`
	Status ShareKubeStatus `json:"status,omitempty"`
}

// DeepCopyInto copies all properties of this object into another object of the same type that is provided as a pointer.
func (in *ShareKube) DeepCopyInto(out *ShareKube) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy creates a new instance of this structure, and then copies the values from the original.
func (in *ShareKube) DeepCopy() *ShareKube {
	if in == nil {
		return nil
	}
	out := new(ShareKube)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject implements the runtime.Object interface.
func (in *ShareKube) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

// DeepCopyInto for ShareKubeSpec
func (in *ShareKubeSpec) DeepCopyInto(out *ShareKubeSpec) {
	*out = *in
	
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = make([]Resource, len(*in))
		copy(*out, *in)
	}
	
	if in.TransformationRules != nil {
		in, out := &in.TransformationRules, &out.TransformationRules
		*out = make([]TransformationRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	
	if in.TargetCluster != nil {
		in, out := &in.TargetCluster, &out.TargetCluster
		*out = new(TargetCluster)
		**out = **in
	}
}

// DeepCopyInto for TransformationRule
func (in *TransformationRule) DeepCopyInto(out *TransformationRule) {
	*out = *in
	if in.RemoveFields != nil {
		in, out := &in.RemoveFields, &out.RemoveFields
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopyInto for ShareKubeStatus
func (in *ShareKubeStatus) DeepCopyInto(out *ShareKubeStatus) {
	*out = *in
	
	if in.CreationTime != nil {
		in, out := &in.CreationTime, &out.CreationTime
		*out = (*in).DeepCopy()
	}
	
	if in.ExpirationTime != nil {
		in, out := &in.ExpirationTime, &out.ExpirationTime
		*out = (*in).DeepCopy()
	}
	
	if in.CopiedResources != nil {
		in, out := &in.CopiedResources, &out.CopiedResources
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

//+kubebuilder:object:root=true

// ShareKubeList contains a list of ShareKube
type ShareKubeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ShareKube `json:"items"`
}

// DeepCopyInto copies all properties of this object into another object of the same type that is provided as a pointer.
func (in *ShareKubeList) DeepCopyInto(out *ShareKubeList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ShareKube, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy creates a new instance of this structure, and then copies the values from the original.
func (in *ShareKubeList) DeepCopy() *ShareKubeList {
	if in == nil {
		return nil
	}
	out := new(ShareKubeList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject implements the runtime.Object interface.
func (in *ShareKubeList) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

func init() {
	SchemeBuilder.Register(&ShareKube{}, &ShareKubeList{})
} 