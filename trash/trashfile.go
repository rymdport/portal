package trash

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// TrashResult is the status for sending a file to the trashcan.
type TrashResult = uint8

const (
	// TrashingFailed indicates that the file could not be sent to the trashcan.
	TrashingFailed TrashResult = 0

	// TrashinSucceeded indicates that the file was sent to the trashcan.
	TrashingSucceeded TrashResult = 1
)

// TrashFile sends a file to the trashcan. Applications are allowed to trash a file if they can open it in r/w mode.
func TrashFile(fd uintptr) (TrashResult, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return 0, err
	}

	data := map[string]dbus.Variant{}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(trashCallName, 0, dbus.UnixFDIndex(fd), data)
	if call.Err != nil {
		return TrashingFailed, call.Err
	}

	var result uint8
	err = call.Store(&result)
	if err != nil {
		return TrashingFailed, err
	}

	return result, nil
}
