package filechooser

import (
	"errors"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const fileChooserCallName = portal.CallBaseName + ".FileChooser"

var errorUnexpectedResponse = errors.New("unexpected response")

func readURIFromResponse(conn *dbus.Conn, call *dbus.Call) ([]string, error) {
	var responsePath dbus.ObjectPath
	err := call.Store(&responsePath)
	if err != nil {
		return nil, err
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(responsePath),
		dbus.WithMatchInterface(portal.RequestInterface),
		dbus.WithMatchMember(portal.ResponseMember),
	)
	if err != nil {
		return nil, err
	}

	dbusChan := make(chan *dbus.Signal)
	conn.Signal(dbusChan)

	response := <-dbusChan
	if len(response.Body) != 2 {
		return nil, errorUnexpectedResponse
	}

	if responseKey, ok := response.Body[0].(uint32); !ok {
		return nil, errorUnexpectedResponse
	} else if responseKey == 1 || responseKey == 2 {
		return nil, nil
	}

	result, ok := response.Body[1].(map[string]dbus.Variant)
	if !ok {
		return nil, errorUnexpectedResponse
	}

	uris, ok := result["uris"].Value().([]string)
	if !ok {
		return nil, errorUnexpectedResponse
	}

	return uris, nil
}
