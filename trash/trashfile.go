package trash

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// TrashFile sends a file to the trashcan. Applications are allowed to trash a file if they can open it in r/w mode.
func TrashFile(fd uintptr) (uint8, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return 0, err
	}

	data := map[string]dbus.Variant{}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(trashCallName, 0, dbus.UnixFDIndex(fd), data)
	if call.Err != nil {
		return 0, err
	}

	var result uint8
	err = call.Store(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
