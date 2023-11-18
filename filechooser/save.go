package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

// SaveOptions contains the options for how files are to be selected.
type SaveOptions struct {
	Modal bool
}

// SaveFile opens a filechooser for selecting where to save a file.
// The chooser will use the supplied title as it's name.
func SaveFile(title string, options *SaveOptions) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	parentWindow := ""
	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(options.Modal),
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".SaveFile", 0, parentWindow, title, data)
	return call.Err
}
