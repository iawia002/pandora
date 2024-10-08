// Code generated by client-gen. DO NOT EDIT.

package v1alpha2

import (
	"context"

	v1alpha2 "github.com/iawia002/pandora/kubernetes/apis/foo/v1alpha2"
	scheme "github.com/iawia002/pandora/kubernetes/generated/clientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// FoosGetter has a method to return a FooInterface.
// A group's client should implement this interface.
type FoosGetter interface {
	Foos() FooInterface
}

// FooInterface has methods to work with Foo resources.
type FooInterface interface {
	Create(ctx context.Context, foo *v1alpha2.Foo, opts v1.CreateOptions) (*v1alpha2.Foo, error)
	Update(ctx context.Context, foo *v1alpha2.Foo, opts v1.UpdateOptions) (*v1alpha2.Foo, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, foo *v1alpha2.Foo, opts v1.UpdateOptions) (*v1alpha2.Foo, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha2.Foo, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha2.FooList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha2.Foo, err error)
	FooExpansion
}

// foos implements FooInterface
type foos struct {
	*gentype.ClientWithList[*v1alpha2.Foo, *v1alpha2.FooList]
}

// newFoos returns a Foos
func newFoos(c *FooV1alpha2Client) *foos {
	return &foos{
		gentype.NewClientWithList[*v1alpha2.Foo, *v1alpha2.FooList](
			"foos",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1alpha2.Foo { return &v1alpha2.Foo{} },
			func() *v1alpha2.FooList { return &v1alpha2.FooList{} }),
	}
}
