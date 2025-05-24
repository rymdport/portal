package screenshot

import (
	"errors"
	"image/color"
	"math"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/request"
)

const pickColorCallName = interfaceName + ".PickColor"

var errInvalidColor = errors.New("invalid color")

// PickerOptions contains options for the color picker.
type PickerOptions struct {
	HandleToken string // A string that will be used as the last element of the handle. Must be a valid object path element.
}

// PickColor obtains the color of a single pixel.
func PickColor(parentWindow string, options *PickerOptions) (*color.RGBA, error) {
	data := map[string]dbus.Variant{}
	if options != nil && options.HandleToken != "" {
		data["handleToken"] = dbus.MakeVariant(options.HandleToken)
	}

	result, err := apis.Call(pickColorCallName, parentWindow, data)
	if err != nil {
		return nil, err
	}

	status, results, err := request.OnSignalResponse(result.(dbus.ObjectPath))
	if err != nil {
		return nil, err
	} else if status == request.Cancelled {
		return nil, nil
	}

	components := results["color"].Value().([]any)
	red := math.Round(math.Max(0.0, math.Min(1.0, components[0].(float64))) * 255)
	green := math.Round(math.Max(0.0, math.Min(1.0, components[1].(float64))) * 255)
	blue := math.Round(math.Max(0.0, math.Min(1.0, components[2].(float64))) * 255)

	return &color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 1,
	}, nil
}
