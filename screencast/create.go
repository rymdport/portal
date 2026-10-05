package screencast

import (
	"github.com/godbus/dbus/v5"

	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/session"
)

const (
	createSessionCallName = interfaceName + ".CreateSession"
	sessionObjectPath     = apis.ObjectPath + "/session/"
)

// CreateSession creates a ScreenCast session.
// A successfully created session can at any time be closed,
// or may at any time be closed by the portal implementation.
func CreateSession() (*Session, error) {
	data := map[string]dbus.Variant{
		"handle_token":         session.GenerateToken(),
		"session_handle_token": session.GenerateToken(),
	}

	result, err := apis.Call(createSessionCallName, data)
	if err != nil {
		return nil, err
	}

	results, err := getResponse(result.(dbus.ObjectPath))
	if err != nil {
		return nil, err
	}

	sessionHandle := dbus.ObjectPath(results["session_handle"].Value().(string))

	return &Session{path: sessionHandle}, nil
}
