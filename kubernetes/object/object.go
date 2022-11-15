package object

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ContainsAnnotation determines whether the object contains an annotation.
func ContainsAnnotation(obj metav1.Object, key string) bool {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		return false
	}
	if _, ok := annotations[key]; ok {
		return true
	}
	return false
}

// GetAnnotation returns the annotation value of the object.
func GetAnnotation(obj metav1.Object, key string) string {
	annotations := obj.GetAnnotations()
	if annotations == nil {
		return ""
	}
	return annotations[key]
}

// ContainsLabel determines whether the object contains a label.
func ContainsLabel(obj metav1.Object, key string) bool {
	labels := obj.GetLabels()
	if labels == nil {
		return false
	}
	if _, ok := labels[key]; ok {
		return true
	}
	return false
}

// GetLabel returns the label value of the object.
func GetLabel(obj metav1.Object, key string) string {
	labels := obj.GetLabels()
	if labels == nil {
		return ""
	}
	return labels[key]
}
