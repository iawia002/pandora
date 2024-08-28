package webhook

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	foov1alpha2 "github.com/iawia002/pandora/kubernetes/apis/foo/v1alpha2"
)

// SetupFooWebhookWithManager ...
func SetupFooWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&foov1alpha2.Foo{}).
		WithDefaulter(&fooDefaulter{}).
		WithValidator(&fooValidator{}).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-foo-example-io-v1alpha2-foo,mutating=true,failurePolicy=fail,sideEffects=None,groups=foo.example.io,resources=foos,verbs=create;update,versions=v1alpha1,name=foo.example.io,admissionReviewVersions=v1
// +kubebuilder:webhook:path=/validate-foo-example-io-v1alpha2-foo,mutating=false,failurePolicy=fail,sideEffects=None,groups=foo.example.io,resources=foos,verbs=create;update,versions=v1alpha1,name=foo.example.io,admissionReviewVersions=v1

type fooDefaulter struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (f *fooDefaulter) Default(_ context.Context, obj runtime.Object) error {
	foo, ok := obj.(*foov1alpha2.Foo)
	if !ok {
		return fmt.Errorf("expected a Foo but got a %T", obj)
	}

	if foo.Spec.Type == "" {
		foo.Spec.Type = foov1alpha2.FooTypeA
	}
	return nil
}

type fooValidator struct{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (f *fooValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return f.validateFoo(ctx, obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (f *fooValidator) ValidateUpdate(ctx context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	return f.validateFoo(ctx, newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (f *fooValidator) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	// TODO: fill in your validation logic upon object deletion.
	return nil, nil
}

func (f *fooValidator) validateFoo(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	foo, ok := obj.(*foov1alpha2.Foo)
	if !ok {
		return nil, fmt.Errorf("expected a Foo but got a %T", obj)
	}

	var allErrs field.ErrorList
	if foo.Spec.Type == "" && foo.Spec.Key == "" {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec", "type", "key"),
			foo.Spec.Type,
			"can't be empty"),
		)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: foov1alpha2.GroupName, Kind: "Foo"}, foo.Name, allErrs)
}
