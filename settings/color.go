package settings

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const readOneCallPath = settingsCallPath + ".ReadOne"

// ColorScheme is the type of color scheme preference that the user has set.
type ColorScheme uint8

const (
	NoPreference = ColorScheme(iota) // Indicates that no appearance preference was set.
	Dark                             // Indicates that dark mode is preferred.
	Light                            // Indicates that light mode is preferred.
)

// GetColorScheme returns the currently set color scheme.
func GetColorScheme() (ColorScheme, error) {
	dbusConn, err := dbus.SessionBus()
	if err != nil {
		return NoPreference, err
	}

	dbusObj := dbusConn.Object(portal.ObjectName, portal.ObjectPath)
	call := dbusObj.Call(
		readOneCallPath,
		dbus.FlagNoAutoStart,
		"org.freedesktop.appearance",
		"color-scheme",
	)
	if call.Err != nil {
		return NoPreference, err
	}

	var value uint8
	if err = call.Store(&value); err != nil {
		return NoPreference, err
	}

	return ColorScheme(value), nil
}
