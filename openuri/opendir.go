package openuri

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const openDirCallName = openURIBaseName + ".OpenDirectory"

// OpenDirectory asks to open the directory containing a local file in the file browser.
func OpenDirectory(parentWindow string, fd uintptr) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	data := map[string]dbus.Variant{}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(openDirCallName, 0, parentWindow, dbus.UnixFDIndex(fd), data)
	return call.Err
}
