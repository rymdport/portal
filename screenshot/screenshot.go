// Package screenshot lets sandboxed applications request a screenshot.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Screenshot.html.
package screenshot

import (
	"github.com/godbus/dbus/v5"

	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const (
	interfaceName      = apis.CallBaseName + ".Screenshot"
	screenshotCallName = interfaceName + ".Screenshot"
)

// ScreenshotOptions represents options for taking a screenshot.
type ScreenshotOptions struct {
	HandleToken string // A string that will be used as the last element of the handle. Must be a valid object path element.
	NotModal    bool   // Whether the dialog should not be modal. Default is no.
	Interactive bool   // Hint whether the dialog should offer customization before taking a screenshot. Default is no. Since version 2.
}

// Screenshot takes a screenshot, and returns the result path as a string.
func Screenshot(parentWindow string, options *ScreenshotOptions) (string, error) {
	data := map[string]dbus.Variant{}

	if options != nil {
		data = map[string]dbus.Variant{
			"modal":       convert.FromBool(!options.NotModal),
			"interactive": convert.FromBool(options.Interactive),
		}

		if options.HandleToken != "" {
			data["handle_token"] = convert.FromString(options.HandleToken)
		}
	}

	result, err := apis.Call(screenshotCallName, parentWindow, data)
	if err != nil {
		return "", err
	}

	return readURIFromResponse(result.(dbus.ObjectPath))
}

func readURIFromResponse(path dbus.ObjectPath) (string, error) {
	status, results, err := request.OnSignalResponse(path)
	if err != nil {
		return "", err
	} else if status >= request.Cancelled {
		return "", nil
	}

	uri := results["uri"].Value().(string)
	return uri, nil
}
