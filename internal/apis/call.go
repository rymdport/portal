package apis

import "github.com/godbus/dbus/v5"

// Call calls the given call name for a portal using passed arguments and returns the output.
func Call(callName string, args ...any) (any, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	obj := conn.Object(ObjectName, ObjectPath)
	call := obj.Call(callName, 0, args...)
	if call.Err != nil {
		return nil, call.Err
	}

	var result any
	err = call.Store(&result)
	return result, err
}
