package settings

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const settingsCallPath = apis.CallBaseName + ".Settings"

// WatchSettingsChange allows setting a function to run each time the portal settings change.
func WatchSettingsChange(callback func(value []any)) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	if err := conn.AddMatchSignal(
		dbus.WithMatchObjectPath(apis.ObjectPath),
		dbus.WithMatchInterface(settingsCallPath),
		dbus.WithMatchMember("SettingChanged"),
	); err != nil {
		return err
	}
	defer conn.Close()

	dbusChan := make(chan *dbus.Signal)
	conn.Signal(dbusChan)

	for sig := range dbusChan {
		callback(sig.Body)
	}

	return nil
}
