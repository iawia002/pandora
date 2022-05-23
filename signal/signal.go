package signal

import (
	"context"
	"time"

	"k8s.io/klog/v2"
)

// After returns a Context that closes after the given duration.
func After(duration time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(duration)
		klog.Info("timeout, closing the context")
		cancel()
	}()
	return ctx
}
