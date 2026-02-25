package wallpaper

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

// SetWallpaperFile sets wallpaper specified as a local file.
func SetWallpaperFile(parentWindow string, fd uintptr, options *SetWallpaperOptions) error {
	unixFD, err := convert.UintptrToUnixFD(fd)
	if err != nil {
		return err
	}

	data := dbusDataFromOptions(options)
	result, err := apis.Call(setWallpaperFileCallName, parentWindow, unixFD, data)
	if err != nil {
		return err
	}

	return readStatusFromResponse(result.(dbus.ObjectPath))
}
