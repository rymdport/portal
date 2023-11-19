package filechooser

import (
	"errors"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const fileChooserCallName = portal.CallBaseName + ".FileChooser"

var errorUnexpectedResponse = errors.New("unexpected response")

func readURIFromResponse(conn *dbus.Conn, call *dbus.Call) ([]string, error) {
	var responcepath dbus.ObjectPath
	err := call.Store(&responcepath)
	if err != nil {
		return nil, err
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(responcepath),
		dbus.WithMatchInterface(portal.RequestInterface),
		dbus.WithMatchMember(portal.ResponseMember),
	)
	if err != nil {
		return nil, err
	}

	dbusChan := make(chan *dbus.Signal)
	conn.Signal(dbusChan)

	responce := <-dbusChan
	if len(responce.Body) != 2 {
		return nil, errorUnexpectedResponse
	}

	result, ok := responce.Body[1].(map[string]dbus.Variant)
	if !ok {
		return nil, errorUnexpectedResponse
	}

	uris, ok := result["uris"].Value().([]string)
	if !ok {
		return nil, errorUnexpectedResponse
	}

	return uris, nil
}
