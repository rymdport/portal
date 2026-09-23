// Package session is the Session interface shared by portal interfaces that
// create long-lived sessions (location, usb, ...).
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

// Close closes the portal session at path.
func Close(path dbus.ObjectPath) error {
	return apis.CallOnObject(path, closeCallName)
}

// GenerateToken returns a random token prefixed with "rymdportal".
func GenerateToken() dbus.Variant {
	str := strings.Builder{}
	str.WriteString("rymdportal")
	a, _ := rand.Int(rand.Reader, big.NewInt(1<<16))
	str.WriteString(strconv.FormatUint(a.Uint64(), 16))
	return convert.FromString(str.String())
}

// OnSignalClosed blocks until the session emits Closed or done is closed.
// Returns (nil, nil) on cancellation. Pass a nil done to wait forever.
func OnSignalClosed(path dbus.ObjectPath, done <-chan struct{}) (map[string]dbus.Variant, error) {
	signal, cleanup, err := apis.ListenOnSignalAt(path, interfaceName, closedMember)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	var response *dbus.Signal
	select {
	case <-done:
		return nil, nil
	case r, ok := <-signal:
		if !ok {
			return nil, portal.ErrUnexpectedResponse
		}
		response = r
	}

	if len(response.Body) != 1 {
		return nil, portal.ErrUnexpectedResponse
	}
	details, ok := response.Body[0].(map[string]dbus.Variant)
	if !ok {
		return nil, portal.ErrUnexpectedResponse
	}
	return details, nil
}
