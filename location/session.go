package location

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
	"github.com/rymdport/portal/internal/session"
)

const (
	startCallName         = interfaceName + ".Start"
	locationUpdatedMember = "LocationUpdated"
)

// StartOptions represents options used when starting a location session.
type StartOptions struct {
	HandleToken string
}

// Session is the value for a created location session.
// The zero value is not usable.
type Session struct {
	path dbus.ObjectPath
}

// Close closes the current session.
func (s *Session) Close() error {
	return session.Close(s.path)
}

// SetOnClosed sets a callback to run when the session is closed by the portal.
func (s *Session) SetOnClosed(callback func()) {
	go func() {
		session.OnSignalClosed(s.path)
		callback()
	}()
}

// SetOnLocationUpdated sets a callback to run when the location changes.
func (s *Session) SetOnLocationUpdated(callback func(Location)) error {
	signal, err := apis.ListenOnSignal(interfaceName, locationUpdatedMember)
	if err != nil {
		return err
	}

	go func() {
		for trigger := range signal {
			if len(trigger.Body) != 2 {
				continue
			}

			if path, ok := trigger.Body[0].(dbus.ObjectPath); !ok || path != s.path {
				continue
			}

			location := trigger.Body[1].(map[string]dbus.Variant)
			timestamp := location["Timestamp"].Value().([]any)
			callback(Location{
				Latitude:  location["Latitude"].Value().(float64),
				Longitude: location["Longitude"].Value().(float64),
				Altitude:  location["Altitude"].Value().(float64),
				Accuracy:  location["Accuracy"].Value().(float64),
				Speed:     location["Speed"].Value().(float64),
				Heading:   location["Heading"].Value().(float64),
				Timestamp: [2]uint64{timestamp[0].(uint64), timestamp[1].(uint64)},
			})
		}
	}()
	return nil
}

// Start the location session. An application can only attempt start a session once.
func (s *Session) Start(parentWindow string, options *StartOptions) error {
	data := map[string]dbus.Variant{}
	if options != nil {
		data["HandleToken"] = convert.FromString(options.HandleToken)
	}

	result, err := apis.Call(startCallName, s.path, parentWindow, data)
	if err != nil {
		return err
	}

	path := result.(dbus.ObjectPath)
	_, _, err = request.OnSignalResponse(path)
	return err
}
