package background

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

// StatusOptions are the options passed into SetStatus.
type StatusOptions struct {
	Message string // A string that will be used as the status message of the application. This should be a single line that can be presented to the user in a list, not a full sentence or paragraph. Must be shorter than 96 characters.
}

// SetStatus sets the status of the application running in background.
func SetStatus(parentWindow string, options StatusOptions) error {
	data := map[string]dbus.Variant{"message": convert.FromString(options.Message)}

	return apis.CallWithoutResult(requestCallName, parentWindow, data)
}
