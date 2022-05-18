package controller

import (
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

// worker returns a long-running function that will dequeue item from the given queue and process it using
// the given function, will re-queue it if an error happens when processing it, otherwise marks it as done.
// maxRetry = 0 means no limit.
func worker(queue workqueue.RateLimitingInterface, syncer func(string) error, duration time.Duration, maxRetry int) func() {
	workFunc := func() bool {
		key, quit := queue.Get()
		if quit {
			return true
		}
		defer queue.Done(key)
		err := syncer(key.(string))
		if err == nil {
			queue.Forget(key)
			return false
		}

		if maxRetry < 1 || queue.NumRequeues(key) <= maxRetry {
			if duration > 0 {
				queue.AddAfter(key, duration)
			} else {
				queue.AddRateLimited(key)
			}
			return false
		}

		queue.Forget(key)
		runtime.HandleError(err)
		klog.Infof("obj %s exceeds max retry limit %d", key, maxRetry)
		return false
	}

	return func() {
		for {
			if quit := workFunc(); quit {
				return
			}
		}
	}
}

// RateLimitedWorker returns a long-running worker function that will re-queue the item by AddRateLimited function
// if an error happens when processing it.
func RateLimitedWorker(queue workqueue.RateLimitingInterface, syncer func(string) error, maxRetry int) func() {
	return worker(queue, syncer, 0, maxRetry)
}

// DelayingWorker returns a long-running worker function that will re-queue the item after the given duration has
// passed if an error happens when processing it.
func DelayingWorker(queue workqueue.RateLimitingInterface, syncer func(string) error, duration time.Duration, maxRetry int) func() {
	return worker(queue, syncer, duration, maxRetry)
}
