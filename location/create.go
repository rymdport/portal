package location

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/session"
)

const createSessionCallName = interfaceName + ".CreateSession"

// Accuracy specifies the requested accuracy.
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
	DistanceThreshold uint32   // Distance threshold in meters. Default is 0.
	TimeThreshold     uint32   // Time threshold in seconds. Default is 0.
	Accuracy          Accuracy // Requested accuracy. Default is EXACT.
}

// CreateSession creates a location session.
// A successfully created session can at any time be closed,
// or may at any time be closed by the portal implementation.
func CreateSession(options *SessionOptions) (*Session, error) {
	data := map[string]dbus.Variant{"session_handle_token": session.GenerateToken()}
	if options != nil {
		if options.Accuracy == unspecified || options.Accuracy > Exact {
			options.Accuracy = Exact
		}
		data["distance-threshold"] = convert.FromUint32(options.DistanceThreshold)
		data["time-threshold"] = convert.FromUint32(options.TimeThreshold)
		data["accuracy"] = convert.FromUint32(options.Accuracy - 1)
	}

	result, err := apis.Call(createSessionCallName, data)
	if err != nil {
		return nil, err
	}

	return &Session{path: result.(dbus.ObjectPath)}, nil
}
