package memorymonitor

import (
	"github.com/godbus/dbus/v5"
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
func OnSignalLowMemoryWarning(callback func(warning LowMemoryWarning)) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(apis.ObjectPath),
		dbus.WithMatchInterface(memorymonitorBaseName),
		dbus.WithMatchMember("LowMemoryWarning"),
	); err != nil {
		return err
	}

	dbusChan := make(chan *dbus.Signal)
	conn.Signal(dbusChan)

	for sig := range dbusChan {
		if len(sig.Body) == 0 {
			continue
		}

		level, ok := sig.Body[0].(byte)
		if !ok {
			continue
		}

		callback(LowMemoryWarning{Level: level})
	}

	return nil
}
