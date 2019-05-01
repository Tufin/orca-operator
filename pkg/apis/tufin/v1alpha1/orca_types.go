package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OrcaSpec defines the desired state of Orca
// +k8s:openapi-gen=true
type OrcaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Namespace          string            `json:"namespace,omitempty"`
	Domain             string            `json:"domain,omitempty"`
	Project            string            `json:"project,omitempty"`
	IngnoredConfigMaps []string          `json:"ignored_config_maps,omitempty"`
	Components         map[string]bool   `json:"components,omitempty"`
	EndPoints          map[string]string `json:"endpoints,omitempty"`
	Images             map[string]string `json:"images,omitempty"`

	// +kubebuilder:validation:Enum=OpenShift,DockerEE,GKE,AKS,EKS,Unknown
	KubePlatform string `json:"kube_platform,omitempty"`
}

// OrcaStatus defines the observed state of Orca
// +k8s:openapi-gen=true
type OrcaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// +kubebuilder:validation:Enum=Ready,Creating,Failed,Unknown,Updated
	Ready string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Orca is the Schema for the orcas API
// +k8s:openapi-gen=true
type Orca struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrcaSpec   `json:"spec,omitempty"`
	Status OrcaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OrcaList contains a list of Orca
type OrcaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items []Orca    `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Orca{}, &OrcaList{})
}
