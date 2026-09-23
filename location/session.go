package location

import (
	"context"

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

// Session is a location session. The zero value is not usable; obtain one
// from CreateSession.
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

// SetOnLocationUpdated sets a callback to run when the location changes.
// The listener runs until Close is called.
func (s *Session) SetOnLocationUpdated(callback func(Location)) error {
	signal, cleanup, err := apis.ListenOnSignalAt(s.path, interfaceName, locationUpdatedMember)
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

func (s *Session) ensureDone() chan struct{} {
	if s.done == nil {
		s.done = make(chan struct{})
	}
	return s.done
}

// Start the location session. An application can only start a session once.
func (s *Session) Start(parentWindow string, options *StartOptions) error {
	return s.StartContext(context.Background(), parentWindow, options)
}

// StartContext is Start with a context.
func (s *Session) StartContext(ctx context.Context, parentWindow string, options *StartOptions) error {
	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}
	_, err := request.SendRequest(ctx, userToken, startCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		return []any{s.path, parentWindow, data}
	})
	return err
}
