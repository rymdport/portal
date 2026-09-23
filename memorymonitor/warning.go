package memorymonitor

import (
	"context"

	"github.com/rymdport/portal/internal/apis"
)

// LowMemoryWarning contains the information given in such a warning.
type LowMemoryWarning struct {
	Level byte // Representing the level of low memory warning.
}

// OnSignalLowMemoryWarning listens for the LowMemoryWarning signal.
// Signal is emitted when a particular low memory situation happens,
// with 0 being the lowest level of memory availability warning,
// and 255 being the highest.
//
// Deprecated: Use OnSignalLowMemoryWarningContext, which can be stopped.
func OnSignalLowMemoryWarning(callback func(warning LowMemoryWarning)) error {
	return OnSignalLowMemoryWarningContext(context.Background(), callback)
}

// OnSignalLowMemoryWarningContext is OnSignalLowMemoryWarning with a context.
// It returns ctx.Err() once ctx is done.
func OnSignalLowMemoryWarningContext(ctx context.Context, callback func(warning LowMemoryWarning)) error {
	signal, cleanup, err := apis.ListenOnSignal(interfaceName, "LowMemoryWarning")
	if err != nil {
		return err
	}
	defer cleanup()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-signal:
			if len(sig.Body) == 0 {
				continue
			}
			level, ok := sig.Body[0].(byte)
			if !ok {
				continue
			}
			callback(LowMemoryWarning{Level: level})
		}
	}
}
