package openuri

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const openFileCallName = openURIBaseName + ".OpenFile"

// OpenFile asks to open a local file.
func OpenFile(parentWindow string, fd uintptr) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	data := map[string]dbus.Variant{}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(openFileCallName, 0, parentWindow, dbus.UnixFD(fd), data)
	return call.Err
}
