package openuri

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const openDirCallName = interfaceName + ".OpenDirectory"

// OpenDirOptions holds optional settings that can be passed to the OpenDir call.
//
// Deprecated: Use [Options] instead.
type OpenDirOptions = Options

// OpenDirectory asks to open the directory containing a local file in the file browser.
// The input parameter fd should be a file descriptor like the one given from [*os.File.Fd] for example.
func OpenDirectory(parentWindow string, fd uintptr, options *Options) error {
	data := readDataFromOptions(options)
	return apis.CallWithoutResult(openDirCallName, parentWindow, dbus.UnixFD(fd), data)
}
