package samplecontroller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/iawia002/lia/kubernetes/controller"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/iawia002/pandora/kubernetes/scheme"
)

const controllerName = "sample-controller"

type nodeController struct {
	nodeLister corelisters.NodeLister
	nodeSynced cache.InformerSynced
	queue      workqueue.RateLimitingInterface
	recorder   record.EventRecorder
}

// NewController returns a new sample controller.
func NewController(kubeClient kubernetes.Interface, nodeInformer coreinformers.NodeInformer) (manager.Runnable, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})

	c := &nodeController{
		nodeLister: nodeInformer.Lister(),
		nodeSynced: nodeInformer.Informer().HasSynced,
		queue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), controllerName),
		recorder:   eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerName}),
	}

	if _, err := nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: c.enqueue,
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldNode := oldObj.(*corev1.Node)
			newNode := newObj.(*corev1.Node)
			if oldNode.ResourceVersion == newNode.ResourceVersion {
				return
			}
			c.enqueue(newObj)
		},
		DeleteFunc: c.enqueue,
	}); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *nodeController) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.nodeSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.Until(controller.RateLimitedWorker(c.queue, c.syncHandler, 5), time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// Check if our nodeController implements necessary interfaces.
var _ manager.Runnable = &nodeController{}
var _ manager.LeaderElectionRunnable = &nodeController{}

// NeedLeaderElection implements the LeaderElectionRunnable interface,
// controllers need to be run in leader election mode.
func (c *nodeController) NeedLeaderElection() bool {
	return true
}

func (c *nodeController) Start(ctx context.Context) error {
	return c.Run(2, ctx.Done())
}

func (c *nodeController) enqueue(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	// custom filter operation
	if !strings.Contains(key, "cd") {
		return
	}
	c.queue.Add(key)
}

func (c *nodeController) syncHandler(key string) error {
	node, err := c.nodeLister.Get(key)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	klog.Info(node.Name)
	// c.recorder.Event(node, corev1.EventTypeNormal, "reason", "message")
	return nil
}
