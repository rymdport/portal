package wallpaper

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// SetWallpaperFile sets wallpaper specified as a local file.
func SetWallpaperFile(parentWindow string, fd uintptr, options *SetWallpaperOptions) error {
	data := dbusDataFromOptions(options)

	result, err := apis.Call(setWallpaperFileCallName, parentWindow, dbus.UnixFD(fd), data)
	if err != nil {
		return err
	}

	return readStatusFromResponse(result.(dbus.ObjectPath))
}
