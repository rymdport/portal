package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

// OpenOptions contains the options for how files are to be selected.
type OpenOptions struct {
	Modal       bool
	Multiple    bool
	Directory   bool
	AcceptLabel string
}

// OpenFile opens a filechooser for selecting a file to open.
// The chooser will use the supplied title as it's name.
func OpenFile(parentWindow, title string, options *OpenOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{
		"modal":     dbus.MakeVariant(options.Modal),
		"multiple":  dbus.MakeVariant(options.Multiple),
		"directory": dbus.MakeVariant(options.Directory),
	}

	if options.AcceptLabel != "" {
		data["accept_label"] = dbus.MakeVariant(options.AcceptLabel)
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".OpenFile", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}
