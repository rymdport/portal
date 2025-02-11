package trash

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const trashFileCallName = interfaceName + ".TrashFile"

// TrashResult is the status for sending a file to the trashcan.
type TrashResult = uint8

const (
	// TrashingFailed indicates that the file could not be sent to the trashcan.
	TrashingFailed TrashResult = 0

	// TrashinSucceeded indicates that the file was sent to the trashcan.
	TrashingSucceeded TrashResult = 1
)

// TrashFile sends a file to the trashcan. Applications are allowed to trash a file if they can open it in r/w mode.
// The input parameter fd should be a file descriptor like the one given from [*os.File.Fd] for example.
func TrashFile(fd uintptr) (TrashResult, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return 0, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(trashFileCallName, 0, dbus.UnixFD(fd))
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
