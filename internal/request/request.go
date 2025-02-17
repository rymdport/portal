// Package request implements handling for the Request interface shared by all portal interfaces.
package request

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
	"github.com/rymdport/portal/internal/apis"
)

const (
	interfaceName  = "org.freedesktop.portal.Request"
	responseMember = "Response"
	closeCallName  = interfaceName + ".Close"
)

// ResponseStatus contains the status of the response.
type ResponseStatus = uint32

const (
	Success   ResponseStatus = 0 // The request was carried out.
	Cancelled ResponseStatus = 1 // The user cancelled the interaction.
	Ended     ResponseStatus = 2 // The user interaction was ended in some other way.
)

// Close closes the portal request to which this object refers and ends all related user interaction (dialogs, etc).
func Close(path dbus.ObjectPath) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	obj := conn.Object(apis.ObjectName, path)
	call := obj.Call(closeCallName, 0)
	return call.Err
}

// OnSignalResponse takes the given dbus connection and tries to read the response object.
// This only works for dbus calls that have an associated response.
func OnSignalResponse(path dbus.ObjectPath) (ResponseStatus, map[string]dbus.Variant, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return Ended, nil, err
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(interfaceName),
		dbus.WithMatchMember(responseMember),
	)
	if err != nil {
		return Ended, nil, err
	}

	dbusChan := make(chan *dbus.Signal)
	conn.Signal(dbusChan)

	response := <-dbusChan
	if len(response.Body) != 2 {
		return Ended, nil, portal.ErrUnexpectedResponse
	}

	status := response.Body[0].(ResponseStatus)
	results := response.Body[1].(map[string]dbus.Variant)
	return status, results, nil
}
