package controllerruntime

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// NodeReconciler is a Reconciler for the Node object.
type NodeReconciler struct {
	client.Client
}

// Reconcile reconciles the Node object.
func (r *NodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	node := &corev1.Node{}
	if err := r.Get(ctx, req.NamespacedName, node); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	klog.Info(node.Name)
	return ctrl.Result{}, nil
}

// InjectClient is used to inject the client into NodeReconciler.
func (r *NodeReconciler) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

// SetupWithManager setups the NodeReconciler with manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return builder.
		ControllerManagedBy(mgr).
		For( // For is required, Owns and Watches are optional
			&corev1.Node{},
			builder.WithPredicates(
				predicate.NewPredicateFuncs(nodeNameFilter),
				predicate.ResourceVersionChangedPredicate{},
			),
		).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 2,
		}).
		Complete(r)
}

// nodeNameFilter is a custom example filter.
func nodeNameFilter(object client.Object) bool {
	return strings.Contains(object.GetName(), "cd")
}
