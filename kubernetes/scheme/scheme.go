package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	foov1alpha1 "github.com/iawia002/pandora/kubernetes/apis/foo/v1alpha1"
	foov1alpha2 "github.com/iawia002/pandora/kubernetes/apis/foo/v1alpha2"
)

// Scheme contains all types of custom Scheme and kubernetes client-go Scheme.
var Scheme = runtime.NewScheme()

func init() {
	_ = clientgoscheme.AddToScheme(Scheme)
	_ = foov1alpha1.Install(Scheme)
	_ = foov1alpha2.Install(Scheme)
}
