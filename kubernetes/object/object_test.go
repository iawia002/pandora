package object

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestContainsAnnotation(t *testing.T) {
	tests := []struct {
		name   string
		obj    metav1.Object
		key    string
		wanted bool
	}{
		{
			name: "contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "aa",
			wanted: true,
		},
		{
			name: "not contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "bb",
			wanted: false,
		},
		{
			name:   "nil test",
			obj:    &corev1.Node{},
			key:    "aa",
			wanted: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsAnnotation(tt.obj, tt.key); got != tt.wanted {
				t.Errorf("ContainsAnnotation() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestGetAnnotation(t *testing.T) {
	tests := []struct {
		name   string
		obj    metav1.Object
		key    string
		wanted string
	}{
		{
			name: "contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "aa",
			wanted: "bb",
		},
		{
			name: "not contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "bb",
			wanted: "",
		},
		{
			name:   "nil test",
			obj:    &corev1.Node{},
			key:    "aa",
			wanted: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAnnotation(tt.obj, tt.key); got != tt.wanted {
				t.Errorf("GetAnnotation() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestContainsLabel(t *testing.T) {
	tests := []struct {
		name   string
		obj    metav1.Object
		key    string
		wanted bool
	}{
		{
			name: "contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "aa",
			wanted: true,
		},
		{
			name: "not contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "bb",
			wanted: false,
		},
		{
			name:   "nil test",
			obj:    &corev1.Node{},
			key:    "aa",
			wanted: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsLabel(tt.obj, tt.key); got != tt.wanted {
				t.Errorf("ContainsLabel() = %v, want %v", got, tt.wanted)
			}
		})
	}
}

func TestGetLabel(t *testing.T) {
	tests := []struct {
		name   string
		obj    metav1.Object
		key    string
		wanted string
	}{
		{
			name: "contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "aa",
			wanted: "bb",
		},
		{
			name: "not contains test",
			obj: &corev1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"aa": "bb",
					},
				},
			},
			key:    "bb",
			wanted: "",
		},
		{
			name:   "nil test",
			obj:    &corev1.Node{},
			key:    "aa",
			wanted: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLabel(tt.obj, tt.key); got != tt.wanted {
				t.Errorf("GetLabel() = %v, want %v", got, tt.wanted)
			}
		})
	}
}
