package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

// OpenOptions contains the options for how files are to be selected.
type OpenOptions struct {
	Modal     bool
	Multiple  bool
	Directory bool
}

// OpenFile opens a filechooser for selecting a file to open.
// The chooser will use the supplied title as it's name.
func OpenFile(title string, options *OpenOptions) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	parentWindow := ""
	data := map[string]dbus.Variant{
		"modal":     dbus.MakeVariant(options.Modal),
		"multiple":  dbus.MakeVariant(options.Multiple),
		"directory": dbus.MakeVariant(options.Directory),
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".OpenFile", 0, parentWindow, title, data)
	return call.Err
}
