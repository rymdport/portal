package trash

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const trashFileCallName = interfaceName + ".TrashFile"

// TrashResult is the status for sending a file to the trashcan.
type TrashResult = uint32

const (
	// TrashingFailed indicates that the file could not be sent to the trashcan.
	TrashingFailed TrashResult = 0

	// TrashinSucceeded indicates that the file was sent to the trashcan.
	TrashingSucceeded TrashResult = 1
)

// TrashFile sends a file to the trashcan. Applications are allowed to trash a file if they can open it in r/w mode.
// The input parameter fd should be a file descriptor like the one given from [*os.File.Fd] for example.
func TrashFile(fd uintptr) (TrashResult, error) {
	result, err := apis.Call(trashFileCallName, dbus.UnixFD(fd))
	if err != nil {
		return TrashingFailed, err
	}

	return result.(TrashResult), nil
}
