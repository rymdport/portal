package networkmonitor

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// OnSignalChanged calls the passed function when the network configuration changes.
func OnSignalChanged(callback func()) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(apis.ObjectPath),
		dbus.WithMatchInterface(interfaceName),
		dbus.WithMatchMember("changed"),
	); err != nil {
		return err
	}

	dbusChan := make(chan *dbus.Signal)
	conn.Signal(dbusChan)

	for range dbusChan {
		callback()
	}

	return nil
}
