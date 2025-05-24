package openuri

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const openFileCallName = interfaceName + ".OpenFile"

// OpenFileOptions holds optional settings that can be passed to the OpenFile call.
//
// Deprecated: Use [Options] instead.
type OpenFileOptions = Options

// OpenFile asks to open a local file.
// The input parameter fd should be a file descriptor like the one given from [*os.File.Fd] for example.
func OpenFile(parentWindow string, fd uintptr, options *Options) error {
	data := readDataFromOptions(options)
	return apis.CallWithoutResult(openFileCallName, parentWindow, dbus.UnixFD(fd), data)
}
