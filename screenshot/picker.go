package screenshot

import (
	"context"
	"image/color"
	"math"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const pickColorCallName = interfaceName + ".PickColor"

// PickerOptions contains options for the color picker.
type PickerOptions struct {
	HandleToken string // A string that will be used as the last element of the handle. Must be a valid object path element.
}

// PickColor obtains the color of a single pixel.
func PickColor(parentWindow string, options *PickerOptions) (*color.RGBA, error) {
	return PickColorContext(context.Background(), parentWindow, options)
}

// PickColorContext is PickColor with a context.
func PickColorContext(ctx context.Context, parentWindow string, options *PickerOptions) (*color.RGBA, error) {
	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}

	resp, err := request.SendRequest(ctx, userToken, pickColorCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		return []any{parentWindow, data}
	})
	if err != nil {
		return nil, err
	} else if resp.Status == request.Cancelled {
		return nil, nil
	}

	components, ok := resp.Results["color"].Value().([]any)
	if !ok || len(components) < 3 {
		return nil, nil
	}
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
