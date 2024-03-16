package settings

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const readOneCallPath = settingsCallPath + ".ReadOne"

// ReadOne reads a single value which may be any valid DBus type. Returns an error on any unknown namespace or key.
func ReadOne(namespace, key string) (any, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(readOneCallPath, 0, namespace, key)
	if call.Err != nil {
		return nil, call.Err
	}

	var value any
	err = call.Store(&value)
	return value, err
}
