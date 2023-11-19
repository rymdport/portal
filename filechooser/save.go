package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

// SaveSingleOptions contains the options for how a file is saved.
type SaveSingleOptions struct {
	Modal bool
}

// SaveFile opens a filechooser for selecting where to save a file.
// The chooser will use the supplied title as it's name.
func SaveFile(title string, options *SaveSingleOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	parentWindow := ""
	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(options.Modal),
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".SaveFile", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}

// SaveMultipleOptions contains the options for how files are saved.
type SaveMultipleOptions struct {
	Modal bool
}

// SaveFiles opens a filechooser for selecting where to save one or more files.
// The chooser will use the supplied title as it's name.
func SaveFiles(title string, options *SaveMultipleOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	parentWindow := ""
	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(options.Modal),
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".SaveFiles", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}
