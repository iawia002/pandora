package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	foov1alpha1 "github.com/iawia002/pandora/kubernetes/apis/example/apis/foo/v1alpha1"
)

// Scheme contains all types of custom Scheme and kubernetes client-go Scheme.
var Scheme = runtime.NewScheme()

func init() {
	_ = clientgoscheme.AddToScheme(Scheme)
	_ = foov1alpha1.AddToScheme(Scheme)
}
