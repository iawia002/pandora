package signal

import (
	"context"
	"time"
)

// After returns a Context that closes after the given duration.
func After(duration time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(duration)
		cancel()
	}()
	return ctx
}
