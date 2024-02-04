package controllerruntime

import (
	"context"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const controllerName = "sample-controller"

// NodeReconciler is a Reconciler for the Node object.
type NodeReconciler struct {
	client.Client

	recorder record.EventRecorder
}

// Reconcile reconciles the Node object.
func (r *NodeReconciler) Reconcile(_ context.Context, node *corev1.Node) (ctrl.Result, error) {
	klog.Info(node.Name)
	// r.recorder.Event(node, corev1.EventTypeNormal, "reason", "message")
	return ctrl.Result{}, nil
}

// SetupWithManager setups the NodeReconciler with manager.
func (r *NodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.Client = mgr.GetClient()
	r.recorder = mgr.GetEventRecorderFor(controllerName)

	// Register timed tasker
	if err := mgr.Add(r); err != nil {
		return err
	}

	return builder.
		ControllerManagedBy(mgr).
		For( // For is required, Owns and Watches are optional
			&corev1.Node{},
			builder.WithPredicates(
				predicate.NewPredicateFuncs(nodeNameFilter),
				predicate.ResourceVersionChangedPredicate{},
			),
		).
		Watches(
			&corev1.Pod{},
			handler.EnqueueRequestsFromMapFunc(func(ctx context.Context, obj client.Object) []reconcile.Request {
				pod := obj.(*corev1.Pod)
				return []reconcile.Request{
					{
						NamespacedName: types.NamespacedName{
							Name: pod.Spec.NodeName,
						},
					},
				}
			}),
			builder.WithPredicates(
				predicate.ResourceVersionChangedPredicate{},
				predicate.Funcs{
					CreateFunc: func(event event.CreateEvent) bool {
						pod := event.Object.(*corev1.Pod)
						return pod.Namespace == metav1.NamespaceSystem
					},
				},
			),
		).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 2,
		}).
		Named(controllerName).
		Complete(reconcile.AsReconciler[*corev1.Node](r.Client, r))
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

// Start starts the timed task.
func (r *NodeReconciler) Start(ctx context.Context) error {
	go wait.Until(func() {
		klog.Infof("current time: %s", time.Now().String())
	}, time.Second*5, ctx.Done())

	return nil
}
