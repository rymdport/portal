package wallpaper

import (
	"net/url"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// SetWallpaperURI sets wallpaper specified as an URI (not a local file)
func SetWallpaperURI(parentWindow string, uri string, options *SetWallpaperOptions) error {
	// Parse URI as an URL !
	url, err := url.Parse(uri)
	if err != nil {
		return err
	}

	data := dbusDataFromOptions(options)

	result, err := apis.Call(setWallpaperURICallName, parentWindow, url.String(), data)
	if err != nil {
		return err
	}

	return readStatusFromResponse(result.(dbus.ObjectPath))
}
