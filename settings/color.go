package settings

import (
	"errors"
	"image/color"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

// ErrNotSet indicates that the value is not set.
var ErrNotSet = errors.New("not set")

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

	dbusObj := dbusConn.Object(apis.ObjectName, apis.ObjectPath)
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

// GetAccentColor returns the currently set accent color.
// If not set, the ErrorNotSet will be returned.
func GetAccentColor() (color.RGBA, error) {
	dbusConn, err := dbus.SessionBus()
	if err != nil {
		return color.RGBA{}, err
	}

	dbusObj := dbusConn.Object(apis.ObjectName, apis.ObjectPath)
	call := dbusObj.Call(
		readOneCallPath,
		dbus.FlagNoAutoStart,
		"org.freedesktop.appearance",
		"accent-color",
	)
	if call.Err != nil {
		return color.RGBA{}, err
	}

	var value []float64
	if err = call.Store(&value); err != nil {
		return color.RGBA{}, err
	}

	if len(value) != 4 {
		return color.RGBA{}, ErrNotSet
	}

	red := value[0] * 255
	if red < 0 || red > 255 {
		return color.RGBA{}, ErrNotSet
	}

	green := value[1] * 255
	if green < 0 || green > 255 {
		return color.RGBA{}, ErrNotSet
	}

	blue := value[2] * 255
	if blue < 0 || blue > 255 {
		return color.RGBA{}, ErrNotSet
	}

	alpha := value[3] * 255
	if alpha < 0 || alpha > 255 {
		return color.RGBA{}, ErrNotSet
	}

	return color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(alpha),
	}, nil
}

// GetContrast returns the currently set contrast setting.
func GetContrast() (Contrast, error) {
	dbusConn, err := dbus.SessionBus()
	if err != nil {
		return NormalContrast, err
	}

	dbusObj := dbusConn.Object(apis.ObjectName, apis.ObjectPath)
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
