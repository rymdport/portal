package settings

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const readOneCallPath = settingsCallPath + ".ReadOne"

// ColorScheme indicates the system’s preferred color scheme.
type ColorScheme uint8

const (
	NoPreference = ColorScheme(iota) // No preference.
	Dark                             // Prefer dark appearance.
	Light                            // Prefer light appearance.
)

// Contrast indicates the system’s preferred contrast level.
type Contrast uint8

const (
	NormalContrast = Contrast(iota) // No preference (normal contrast)
	HigherContrast                  // Higher contrast
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

	if value > 2 {
		value = 0 // Unknown values should be treated as 0 (no preference).
	}

	return ColorScheme(value), nil
}

// GetContrast returns the currently set contrast setting.
func GetContrast() (Contrast, error) {
	dbusConn, err := dbus.SessionBus()
	if err != nil {
		return NormalContrast, err
	}

	dbusObj := dbusConn.Object(portal.ObjectName, portal.ObjectPath)
	call := dbusObj.Call(
		readOneCallPath,
		dbus.FlagNoAutoStart,
		"org.freedesktop.appearance",
		"contrast",
	)
	if call.Err != nil {
		return NormalContrast, err
	}

	var value uint8
	if err = call.Store(&value); err != nil {
		return NormalContrast, err
	}

	if value > 1 {
		value = 0 // Unknown values should be treated as 0 (no preference).
	}

	return Contrast(value), nil
}
