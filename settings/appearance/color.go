package appearance

import (
	"errors"
	"image/color"
	"math"

	"github.com/rymdport/portal/settings"
)

// ErrNotSet indicates that the value is not set.
var ErrNotSet = errors.New("not set")

// ColorScheme indicates the system’s preferred color scheme.
type ColorScheme uint8

const (
	NoPreference = ColorScheme(iota) // No preference.
	Dark                             // Prefer dark appearance.
	Light                            // Prefer light appearance.
)

// GetColorScheme returns the currently set color scheme.
func GetColorScheme() (ColorScheme, error) {
	value, err := settings.ReadOne(Namespace, "color-scheme")
	if err != nil {
		return NoPreference, err
	}

	return ValueToColorScheme(value)
}

// GetAccentColor returns the currently set accent color.
// If not set, the ErrorNotSet will be returned.
func GetAccentColor() (*color.RGBA, error) {
	value, err := settings.ReadOne(Namespace, "accent-color")
	if err != nil {
		return nil, ErrNotSet
	}

	return ValueToAccentColor(value)
}

// ValueToColorScheme converts a read value to a ColorScheme type.
// This is useful when for example parsing a value from the callback
// in [settings.SignalOnSettingChanged] or a value from [settings.ReadOne].
func ValueToColorScheme(value any) (ColorScheme, error) {
	result, ok := value.(uint32)
	if !ok {
		return NoPreference, ErrNotSet
	}

	if result > 2 {
		return 0, nil // Unknown values should be treated as 0 (no preference).
	}

	return ColorScheme(result), nil
}

// ValueToAccentColor converts a read value to an accent color type.
// This is useful when for example parsing a value from the callback
// in [settings.SignalOnSettingChanged] or a value from [settings.ReadOne].
func ValueToAccentColor(value any) (*color.RGBA, error) {
	result, ok := value.([]float64)
	if !ok || len(result) != 4 {
		return nil, ErrNotSet
	}

	red := math.Round(math.Max(0.0, math.Min(1.0, result[0])) * 255)
	green := math.Round(math.Max(0.0, math.Min(1.0, result[1])) * 255)
	blue := math.Round(math.Max(0.0, math.Min(1.0, result[2])) * 255)
	alpha := math.Round(math.Max(0.0, math.Min(1.0, result[3])) * 255)

	return &color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: uint8(alpha),
	}, nil
}
