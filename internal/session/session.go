// Package session is a shared session interface between all portal interfaces that involve long lived sessions. When a method that creates a session is called, if successful, the reply will include a session handle (i.e. object path) for a Session object, which will stay alive for the duration of the session.
package session

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

const (
	interfaceName = "org.freedesktop.portal.Session"
	closedMember  = "Closed"
	closeCallName = interfaceName + ".Close"
)

// Close closes the portal session to which this object refers and ends all related user interaction (dialogs, etc).
func Close(path dbus.ObjectPath) error {
	return apis.CallOnObject(path, closeCallName)
}

// GenerateToken generates a random token string prefixed with "rymdportal".
func GenerateToken() dbus.Variant {
	str := strings.Builder{}
	str.WriteString("rymdportal")
	a, _ := rand.Int(rand.Reader, big.NewInt(1<<16))
	str.WriteString(strconv.FormatUint(a.Uint64(), 16))
	return convert.FromString(str.String())
}

// OnSignalClosed takes the given dbus connection and listens for the closed signal.
// The signal is emitted when a session is closed.
// The content of details is specified by the interface creating the session.
func OnSignalClosed(path dbus.ObjectPath) (map[string]dbus.Variant, error) {
	signal, err := apis.ListenOnSignalAt(path, interfaceName, closedMember)
	if err != nil {
		return nil, err
	}

	for response := range signal {
		if response.Path != path {
			continue
		}
		if len(response.Body) != 1 {
			return nil, portal.ErrUnexpectedResponse
		}

		details := response.Body[0].(map[string]dbus.Variant)
		return details, nil
	}

	return nil, portal.ErrUnexpectedResponse
}
