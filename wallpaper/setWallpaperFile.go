package wallpaper

import (
	"errors"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// SetWallpaperFile sets wallpaper specified as a local file
func SetWallpaperFile(parentWindow string, file *os.File, options *SetWallpaperOptions) error {
	data := dbusDataFromOptions(options)

	if unixFDSupport, err := checkDbusUnixFDSupport(); err != nil {
		return err
	} else if !unixFDSupport {
		return errors.New("DBus connection does not support passing unix file descriptors")
	}

	result, err := apis.Call(setWallpaperFileCallName, parentWindow, dbus.UnixFD(file.Fd()), data)
	if err != nil {
		return err
	}

	return readStatusFromResponse(result.(dbus.ObjectPath))
}

func checkDbusUnixFDSupport() (bool, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return false, err
	}
	return conn.SupportsUnixFDs(), nil
}
