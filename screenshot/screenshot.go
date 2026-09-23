// Package screenshot lets sandboxed applications request a screenshot.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Screenshot.html.
package screenshot

import (
	"context"

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
	return ScreenshotContext(context.Background(), parentWindow, options)
}

// ScreenshotContext is like [Screenshot] but accepts a [context.Context] that
// can be used to cancel the request.
func ScreenshotContext(ctx context.Context, parentWindow string, options *ScreenshotOptions) (string, error) {
	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}

	resp, err := request.SendRequest(ctx, userToken, screenshotCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		if options != nil {
			data["modal"] = convert.FromBool(!options.NotModal)
			data["interactive"] = convert.FromBool(options.Interactive)
		}
		return []any{parentWindow, data}
	})
	if err != nil {
		return "", err
	} else if resp.Status >= request.Cancelled {
		return "", nil
	}

	uri, _ := resp.Results["uri"].Value().(string)
	return uri, nil
}
