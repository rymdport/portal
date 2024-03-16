package settings

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const readAllCallPath = settingsCallPath + ".ReadAll"

// ReadAll returns all values for the corresponding namespaces passed.
// If namespaces is an empty array or contains an empty string it matches all.
// Globbing is supported but only for trailing sections, e.g. “org.example.*”.
func ReadAll(namespaces []string) (map[string](map[string]dbus.Variant), error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(readAllCallPath, 0, namespaces)
	if call.Err != nil {
		return nil, call.Err
	}

	var value map[string](map[string]dbus.Variant)
	err = call.Store(&value)
	return value, err
}
