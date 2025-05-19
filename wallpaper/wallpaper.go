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

// WallpaperLocation is the type of the parameter SetOn of wallpaper options
type WallpaperLocation uint

const (
	Background WallpaperLocation = iota // Set wallpaper of Background
	Lockscreen                          // Set wallpaper of Locksreen
	Both                                // Set wallpaper of both background and lockscreen
)

var wallpaperLocationName = map[WallpaperLocation]string{
	Background: "background",
	Lockscreen: "lockscreen",
	Both:       "both",
}

func (l WallpaperLocation) String() string {
	return wallpaperLocationName[l]
}

// SetWallpaperOptions contains the options of backgound change
type SetWallpaperOptions struct {
	ShowPreview bool              // Whether to show a preview of the picture. Note that the portal may decide to show a preview even if this option is not set.
	SetOn       WallpaperLocation // Where to set the wallpaper. Possible values are Background, Lockscreen, or Both constants
}

func dbusDataFromOptions(options *SetWallpaperOptions) map[string]dbus.Variant {
	data := map[string]dbus.Variant{}
	if options != nil {
		data["show-preview"] = convert.FromBool(options.ShowPreview)
		data["set-on"] = convert.FromString(wallpaperLocationName[options.SetOn])
	}
	return data
}

func readStatusFromResponse(path dbus.ObjectPath) error {
	status, _, err := request.OnSignalResponse(path)
	if err != nil {
		return err
	}
	switch status {
	case request.Cancelled:
		return errors.New("operation cancelled by User")
	case request.Ended:
		return errors.New("operation cancelled by system")
	case request.Success:
		return nil
	default:
		return errors.New("unknown status code")
	}
	return nil
}
