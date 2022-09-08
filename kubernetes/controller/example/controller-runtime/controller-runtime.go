package controllerruntime

import (
	"context"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
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

// Timed task example

// Check if our NodeReconciler implements necessary interfaces.
var _ manager.Runnable = &NodeReconciler{}
var _ manager.LeaderElectionRunnable = &NodeReconciler{}

// NeedLeaderElection implements the LeaderElectionRunnable interface,
// controllers need to be run in leader election mode.
func (r *NodeReconciler) NeedLeaderElection() bool {
	return true
}

func (r *NodeReconciler) Start(ctx context.Context) error {
	go wait.Until(func() {
		klog.Infof("current time: %s", time.Now().String())
	}, time.Second*5, ctx.Done())

	return nil
}
