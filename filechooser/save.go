package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

// SaveSingleOptions contains the options for how a file is saved.
type SaveSingleOptions struct {
	Modal       bool
	AcceptLabel string
	FileName    string
	Location    string
}

// SaveFile opens a filechooser for selecting where to save a file.
// The chooser will use the supplied title as it's name.
func SaveFile(parentWindow, title string, options *SaveSingleOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(options.Modal),
	}

	if options.AcceptLabel != "" {
		data["accept_label"] = dbus.MakeVariant("") //dbus.MakeVariant(options.AcceptLabel)
	}

	if options.FileName != "" {
		data["current_name"] = dbus.MakeVariant(options.FileName)
	}

	if options.Location != "" {
		nullTerminatedByteString := []byte(options.Location + "\000")
		data["current_folder"] = dbus.MakeVariant(nullTerminatedByteString)
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
	Modal       bool
	AcceptLabel string
	Location    string
}

// SaveFiles opens a filechooser for selecting where to save one or more files.
// The chooser will use the supplied title as it's name.
func SaveFiles(parentWindow, title string, options *SaveMultipleOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(options.Modal),
	}

	if options.AcceptLabel != "" {
		data["accept_label"] = dbus.MakeVariant(options.AcceptLabel)
	}

	if options.Location != "" {
		nullTerminatedByteString := []byte(options.Location + "\000")
		data["current_folder"] = dbus.MakeVariant(nullTerminatedByteString)
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".SaveFiles", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}
