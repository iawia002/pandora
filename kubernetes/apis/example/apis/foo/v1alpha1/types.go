package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	// FooFinalizer ...
	FooFinalizer = "finalizers.example.io/foo-bar"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced
// +kubebuilder:resource:scope=Cluster,shortName="fo"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Key",type=string,JSONPath=".spec.key"

// Foo ...
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FooSpec `json:"spec"`
	// +optional
	Status FooStatus `json:"status,omitempty"`
}

// +kubebuilder:validation:Enum=A;B

// FooType ...
type FooType string

const (
	// FooTypeA ...
	FooTypeA FooType = "A"
	// FooTypeB ...
	FooTypeB FooType = "B"
)

// FooSpec ...
type FooSpec struct {
	// +kubebuilder:default=A
	Type FooType `json:"type"`
	// +kubebuilder:validation:MinLength=2
	Key string `json:"key"`
	// +optional
	Value string `json:"value,omitempty"`
}

// FooStatus ...
type FooStatus struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Raw *runtime.RawExtension `json:"raw,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList ...
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}
