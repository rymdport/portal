package filechooser

import (
	"errors"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

var errorUnexpectedResponce = errors.New("unexpected responce")

// OpenOptions contains the options for how files are to be selected.
type OpenOptions struct {
	Modal     bool
	Multiple  bool
	Directory bool
}

// OpenFile opens a filechooser for selecting a file to open.
// The chooser will use the supplied title as it's name.
func OpenFile(title string, options *OpenOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	parentWindow := ""
	data := map[string]dbus.Variant{
		"modal":     dbus.MakeVariant(options.Modal),
		"multiple":  dbus.MakeVariant(options.Multiple),
		"directory": dbus.MakeVariant(options.Directory),
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".OpenFile", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	var responcepath dbus.ObjectPath
	err = call.Store(&responcepath)
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
		return nil, errorUnexpectedResponce
	}

	result, ok := responce.Body[1].(map[string]dbus.Variant)
	if !ok {
		return nil, errorUnexpectedResponce
	}

	uris, ok := result["uris"].Value().([]string)
	if !ok {
		return nil, errorUnexpectedResponce
	}

	return uris, nil
}
