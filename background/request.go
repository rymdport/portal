package background

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const requestCallName = interfaceName + ".RequestBackground"

// RequestOptions holds the options used for RequestBackground.
type RequestOptions struct {
	HandleToken     string   // A string that will be used as the last element of the handle. Must be a valid object path element.
	Reason          string   // User-visible reason for the request.
	Autostart       bool     // Specify if the app wants to be started automatically at login.
	Commandline     []string // Commandline to use add when autostarting at login. If this is not specified, the Exec key from the desktop file will be used.
	DbusActivatable bool     // Specify if the app should use D-Bus activation for autostart.
}

// RequestResult is the result returned from RequestBackground.
type RequestResult struct {
	Background bool // True if the application is allowed to run in the background.
	Autostart  bool // True if the application will be autostarted.
}

// RequestBackground requests that the application is allowed to run in the background.
func RequestBackground(parentWindow string, options *RequestOptions) (*RequestResult, error) {
	data := map[string]dbus.Variant{}
	if options != nil {
		data["autostart"] = convert.FromBool(options.Autostart)
		data["dbus-activatable"] = convert.FromBool(options.DbusActivatable)

		if options.HandleToken != "" {
			data["handle_token"] = convert.FromString(options.HandleToken)
		}

		if options.Reason != "" {
			data["reason"] = convert.FromString(options.Reason)
		}

		if len(options.Commandline) != 0 {
			data["commandline"] = dbus.MakeVariant(options.Commandline) // TODO: Might want to create fast converter for []string.
		}
	}

	result, err := apis.Call(requestCallName, parentWindow, data)
	if err != nil {
		return nil, err
	}

	status, results, err := request.OnSignalResponse(result.(dbus.ObjectPath))
	if err != nil {
		return nil, err
	} else if status == request.Cancelled {
		return nil, nil // Cancelled by user.
	}

	background := results["background"].Value().(bool)
	autostart := results["autostart"].Value().(bool)
	return &RequestResult{Background: background, Autostart: autostart}, nil
}
