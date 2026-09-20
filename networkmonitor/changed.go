package networkmonitor

import (
	"context"

	"github.com/rymdport/portal/internal/apis"
)

// OnSignalChanged calls the passed function when the network configuration changes.
//
// Deprecated: Use OnSignalChangedContext, which can be stopped.
func OnSignalChanged(callback func()) error {
	return OnSignalChangedContext(context.Background(), callback)
}

// OnSignalChangedContext is OnSignalChanged with a context.
// It returns ctx.Err() once ctx is done.
func OnSignalChangedContext(ctx context.Context, callback func()) error {
	signal, cleanup, err := apis.ListenOnSignal(interfaceName, "changed")
	if err != nil {
		return err
	}
	defer cleanup()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-signal:
			callback()
		}
	}
}
