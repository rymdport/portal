package wallpaper

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// SetWallpaperURI sets wallpaper specified as a URI.
func SetWallpaperURI(parentWindow string, uri string, options *SetWallpaperOptions) error {
	data := dbusDataFromOptions(options)

	result, err := apis.Call(setWallpaperURICallName, parentWindow, uri, data)
	if err != nil {
		return err
	}

	return readStatusFromResponse(result.(dbus.ObjectPath))
}
