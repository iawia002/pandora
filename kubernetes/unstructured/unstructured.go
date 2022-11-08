package unstructured

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ConvertToUnstructured converts a typed object to an unstructured object.
func ConvertToUnstructured(obj interface{}) (*unstructured.Unstructured, error) {
	uncastObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{Object: uncastObj}, nil
}

//nolint:gofmt,goimports
// ConvertToTyped converts an unstructured object to a typed object.
// Usage:
// 	node := &corev1.Node{}
// 	ConvertToTyped(object, node)
func ConvertToTyped(obj *unstructured.Unstructured, typedObj interface{}) error {
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), typedObj); err != nil {
		return err
	}
	return nil
}
