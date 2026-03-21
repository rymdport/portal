package usb

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/session"
)

const deviceEventsMember = "DeviceEvents"

// DeviceEvent represents a USB device event.
// Each element of the events array is composed of the following fields.
type DeviceEvent struct {
	Action string           // Type of event that occurred. One of "add", "change", or "remove".
	ID     string           // Device ID that the event occurred on.
	Device DeviceProperties // Device properties attached to the ID. See EnumerateDevices for a list of all the properties that may be present in the vardict.
}

// Session is the value for a created USB monitoring session.
// The zero value is not usable.
type Session struct {
	path dbus.ObjectPath
}

// Close closes the current session.
func (s *Session) Close() error {
	return session.Close(s.path)
}

// SetOnClosed sets a callback to run when the session is closed by the portal.
func (s *Session) SetOnClosed(callback func(error)) {
	go func() {
		_, err := session.OnSignalClosed(s.path)
		callback(err)
	}()
}

// SetOnDeviceEvents sets a callback to run when one or more USB devices have been added, changed, or removed.
// The DeviceEvents signal is only emitted for active sessions created with CreateSession.
func (s *Session) SetOnDeviceEvents(callback func([]DeviceEvent)) error {
	signal, err := apis.ListenOnSignalAt(s.path, interfaceName, deviceEventsMember)
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

			rawEvents := trigger.Body[1].([][]any)
			events := make([]DeviceEvent, len(rawEvents))
			for i, raw := range rawEvents {
				events[i] = DeviceEvent{
					Action: raw[0].(string),
					ID:     raw[1].(string),
					Device: raw[2].(map[string]dbus.Variant),
				}
			}

			callback(events)
		}
	}()
	return nil
}
