package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate-foo-example-io-v1alpha1-foo,mutating=true,failurePolicy=fail,sideEffects=None,groups=foo.example.io,resources=foos,verbs=create;update,versions=v1alpha1,name=foo.example.io,admissionReviewVersions=v1
// +kubebuilder:webhook:path=/validate-foo-example-io-v1alpha1-foo,mutating=false,failurePolicy=fail,sideEffects=None,groups=foo.example.io,resources=foos,verbs=create;update,versions=v1alpha1,name=foo.example.io,admissionReviewVersions=v1

// SetupWebhookWithManager ...
func (r *Foo) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

var _ webhook.Defaulter = &Foo{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Foo) Default() {
	if r.Spec.Type == "" {
		r.Spec.Type = FooTypeA
	}
}

var _ webhook.Validator = &Foo{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateCreate() (admission.Warnings, error) {
	return r.validateFoo()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	return r.validateFoo()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Foo) ValidateDelete() (admission.Warnings, error) {
	// TODO: fill in your validation logic upon object deletion.
	return nil, nil
}

func (r *Foo) validateFoo() (admission.Warnings, error) {
	var allErrs field.ErrorList
	if r.Spec.Type == "" && r.Spec.Key == "" {
		allErrs = append(allErrs, field.Invalid(
			field.NewPath("spec", "type", "key"),
			r.Spec.Type,
			"can't be empty"),
		)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(schema.GroupKind{Group: GroupName, Kind: "Foo"}, r.Name, allErrs)
}
