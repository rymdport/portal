package usb

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/session"
)

const deviceEventsMember = "DeviceEvents"

// DeviceEvent represents a USB device event.
type DeviceEvent struct {
	Action string           // "add", "change", or "remove"
	ID     string           // device ID
	Device DeviceProperties // device properties; see EnumerateDevices
}

// Session is a USB monitoring session. The zero value is not usable; obtain
// one from CreateSession.
type Session struct {
	path     dbus.ObjectPath
	done     chan struct{} // closed by Close to stop listeners
	cleanups []func()
}

// Close releases listener goroutines and closes the portal session.
func (s *Session) Close() error {
	if s.done != nil {
		select {
		case <-s.done: // already closed
		default:
			close(s.done)
		}
	}
	for _, c := range s.cleanups {
		c()
	}
	s.cleanups = nil
	return session.Close(s.path)
}

// SetOnClosed sets a callback to run when the portal closes the session.
// The goroutine exits on Close even if the signal never arrives.
func (s *Session) SetOnClosed(callback func(error)) {
	done := s.ensureDone()
	go func() {
		_, err := session.OnSignalClosed(s.path, done)
		select {
		case <-done:
			return // cancelled by Close
		default:
		}
		callback(err)
	}()
}

// SetOnDeviceEvents sets a callback to run when one or more USB devices have been added, changed, or removed.
// The DeviceEvents signal is only emitted for active sessions created with CreateSession.
// The listener runs until Close is called.
func (s *Session) SetOnDeviceEvents(callback func([]DeviceEvent)) error {
	signal, cleanup, err := apis.ListenOnSignalAt(s.path, interfaceName, deviceEventsMember)
	if err != nil {
		return err
	}
	s.cleanups = append(s.cleanups, cleanup)
	done := s.ensureDone()

	go func() {
		for {
			var trigger *dbus.Signal
			select {
			case <-done:
				return
			case trigger = <-signal:
			}
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

func (s *Session) ensureDone() chan struct{} {
	if s.done == nil {
		s.done = make(chan struct{})
	}
	return s.done
}
