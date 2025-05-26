// Package wallpaper is a simple interface that lets sandboxed applications set the user’s desktop background picture.
package wallpaper

import (
	"errors"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const (
	interfaceName            = apis.CallBaseName + ".Wallpaper"
	setWallpaperURICallName  = interfaceName + ".SetWallpaperURI"
	setWallpaperFileCallName = interfaceName + ".SetWallpaperFile"
)

// Location specifies where to set the wallpaper.
type Location string

const (
	Background Location = "background" // Set wallpaper of background
	Lockscreen Location = "lockscreen" // Set wallpaper of lockscreen
	Both       Location = "both"       // Set wallpaper of both background and lockscreen
)

// SetWallpaperOptions contains the options of the wallpaper change.
type SetWallpaperOptions struct {
	ShowPreview bool     // Whether to show a preview of the picture. Note that the portal may decide to show a preview even if this option is not set.
	SetOn       Location // Where to set the wallpaper. Possible values are Background, Lockscreen, or Both constants
}

func dbusDataFromOptions(options *SetWallpaperOptions) map[string]dbus.Variant {
	data := map[string]dbus.Variant{}
	if options != nil {
		data["show-preview"] = convert.FromBool(options.ShowPreview)

		if options.SetOn != "" {
			data["set-on"] = convert.FromString(string(options.SetOn))
		}
	}
	return data
}

func readStatusFromResponse(path dbus.ObjectPath) error {
	status, _, err := request.OnSignalResponse(path)
	if err != nil {
		return err
	}

	if status == request.Ended {
		return errors.New("operation cancelled by system")
	}

	return nil
}
