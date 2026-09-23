// Package request implements the Request interface shared by portal methods
// that return a handle and emit a Response signal asynchronously.
package request

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
	"github.com/rymdport/portal/internal/apis"
)

// https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Request.html
const (
	interfaceName  = "org.freedesktop.portal.Request"
	responseMember = "Response"
	closeCallName  = interfaceName + ".Close"

	requestPathPrefix = "/org/freedesktop/portal/desktop/request/"
)

var errNoBusName = errors.New("connection has no unique bus name")

// ResponseStatus of a portal Response signal.
type ResponseStatus = uint32

const (
	Success   ResponseStatus = 0
	Cancelled ResponseStatus = 1
	Ended     ResponseStatus = 2 // closed by the system or other non-user reason
)

// Response is the payload of a portal Request's Response signal.
type Response struct {
	Handle  dbus.ObjectPath
	Status  ResponseStatus
	Results map[string]dbus.Variant
}

// Close closes the portal Request at path.
func Close(path dbus.ObjectPath) error {
	return apis.CallOnObject(path, closeCallName)
}

// closeNoReply asks the portal to close the Request without waiting for a reply.
func closeNoReply(conn *dbus.Conn, path dbus.ObjectPath) {
	conn.Object(apis.ObjectName, path).Call(closeCallName, dbus.FlagNoReplyExpected)
}

func generateToken() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return "rymdportal" + hex.EncodeToString(b[:]), nil
}

// buildRequestPath is the Request path the portal will use for a call
// made by sender with the given handle_token. See the spec link above.
func buildRequestPath(sender, token string) dbus.ObjectPath {
	sender = strings.TrimPrefix(sender, ":")
	sender = strings.ReplaceAll(sender, ".", "_")
	return dbus.ObjectPath(requestPathPrefix + sender + "/" + token)
}

func expectedHandle(conn *dbus.Conn, token string) (dbus.ObjectPath, error) {
	names := conn.Names()
	if len(names) == 0 {
		return "", errNoBusName
	}
	return buildRequestPath(names[0], token), nil
}

// SendRequest dispatches a portal method that returns a Request handle and
// waits for its Response signal. buildArgs receives the handle_token to
// write into the call's options map; pass token="" to generate one.
//
// The Response subscription is installed before the call to avoid missing
// the signal on fast backends. Cancelling ctx asks the portal to close the
// Request without waiting for a reply.
func SendRequest(
	ctx context.Context,
	token, callName string,
	buildArgs func(token string) []any,
) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	conn, err := dbus.SessionBus()
	if err != nil {
		return Response{}, err
	}

	if token == "" {
		if token, err = generateToken(); err != nil {
			return Response{}, err
		}
	}
	expected, err := expectedHandle(conn, token)
	if err != nil {
		return Response{}, err
	}

	signal, cleanup, err := apis.ListenOnSignalAt(expected, interfaceName, responseMember)
	if err != nil {
		return Response{}, err
	}
	defer cleanup()

	call := conn.Object(apis.ObjectName, apis.ObjectPath).Call(callName, 0, buildArgs(token)...)
	if call.Err != nil {
		return Response{}, call.Err
	}

	var handle dbus.ObjectPath
	if err := call.Store(&handle); err != nil {
		return Response{}, err
	}
	if handle != expected {
		closeNoReply(conn, handle)
		return Response{Handle: handle}, fmt.Errorf("portal returned request path %q, expected %q", handle, expected)
	}

	var sig *dbus.Signal
	select {
	case <-ctx.Done():
		closeNoReply(conn, handle)
		return Response{Handle: handle}, ctx.Err()
	case sig = <-signal:
	}

	if len(sig.Body) != 2 {
		return Response{Handle: handle}, portal.ErrUnexpectedResponse
	}
	status, ok := sig.Body[0].(ResponseStatus)
	if !ok {
		return Response{Handle: handle}, portal.ErrUnexpectedResponse
	}
	results, ok := sig.Body[1].(map[string]dbus.Variant)
	if !ok {
		return Response{Handle: handle}, portal.ErrUnexpectedResponse
	}
	return Response{Handle: handle, Status: status, Results: results}, nil
}
