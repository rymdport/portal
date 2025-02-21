package location

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

const createSessionCallName = interfaceName + ".CreateSession"

// Accuracy specifies the requested location accuracy.
type Accuracy = uint32

const (
	unspecified Accuracy = iota // Shift values up by one so we know if the field if unspecified.
	None
	Country
	City
	Neighborhood
	Street
	Exact
)

// SessionOptions contains options used when creating the location session.
type SessionOptions struct {
	SessionHandleToken string   // A string that will be used as the last element of the session handle. Must be a valid object path element.
	DistanceThreshold  uint32   // Distance threshold in meters. Default is 0.
	TimeThreshold      uint32   // Time threshold in seconds. Default is 0.
	Accuracy           Accuracy // Requested accuracy. Default is EXACT.
}

// CreateSession creates a location session.
// A successfully created session can at any time be closed,
// or may at any time be closed by the portal implementation.
func CreateSession(options SessionOptions) (dbus.ObjectPath, error) {
	if options.Accuracy == unspecified {
		options.Accuracy = Exact
	}

	data := map[string]dbus.Variant{
		"session_handle_token": convert.FromString(options.SessionHandleToken),
		"distance-threshold":   convert.FromUint32(options.DistanceThreshold),
		"time-threshold":       convert.FromUint32(options.TimeThreshold),
		"accuracy":             convert.FromUint32(options.Accuracy - 1),
	}

	result, err := apis.Call(createSessionCallName, data)
	return result.(dbus.ObjectPath), err
}
