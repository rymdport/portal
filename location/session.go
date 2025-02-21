package location

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/request"
	"github.com/rymdport/portal/internal/session"
)

const (
	startCallName         = interfaceName + ".Start"
	locationUpdatedMember = "LocationUpdated"
)

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

			handle := trigger.Body[0].(dbus.ObjectPath)
			location := trigger.Body[1].(map[string]any)
			if handle != s.path {
				continue
			}

			callback(Location{
				Latitude:  location["Latitude"].(float64),
				Longitude: location["Latitude"].(float64),
				Altitude:  location["Altitude"].(float64),
				Accuracy:  location["Accuracy"].(float64),
				Speed:     location["Speed"].(float64),
				Heading:   location["Heading"].(float64),
				Timestamp: location["Timestamp"].([2]uint64),
			})
		}
	}()
	return nil
}

// Start the location session. An application can only attempt start a session once.
func (s *Session) Start(parentWindow string) error {
	return start(s.path, parentWindow)
}

// start starts the location session. An application can only attempt start a session once.
func start(sessionHandle dbus.ObjectPath, parentWindow string) error {
	data := map[string]dbus.Variant{}
	result, err := apis.Call(startCallName, sessionHandle, parentWindow, data)
	if err != nil {
		return err
	}

	path := result.(dbus.ObjectPath)
	status, results, err := request.OnSignalResponse(path)
	fmt.Println(status, results)
	return err
}
