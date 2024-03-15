package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

// SaveFileOptions contains the options for how a file is saved.
type SaveFileOptions struct {
	AcceptLabel   string // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool   // Whether the dialog should not be modal.
	CurrentName   string // Suggested name of the file.
	CurrentFolder string // Suggested folder in which the file should be saved.
}

// SaveFile opens a filechooser for selecting where to save a file.
// The chooser will use the supplied title as it's name.
func SaveFile(parentWindow, title string, options *SaveFileOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(!options.NotModal),
	}

	if options.AcceptLabel != "" {
		data["accept_label"] = dbus.MakeVariant("") // dbus.MakeVariant(options.AcceptLabel)
	}

	if options.CurrentName != "" {
		data["current_name"] = dbus.MakeVariant(options.CurrentName)
	}

	if options.CurrentFolder != "" {
		nullTerminatedByteString := []byte(options.CurrentFolder + "\000")
		data["current_folder"] = dbus.MakeVariant(nullTerminatedByteString)
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".SaveFile", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}

// SaveFilesOptions contains the options for how files are saved.
type SaveFilesOptions struct {
	AcceptLabel   string // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool   // Whether the dialog should be modal.
	CurrentFolder string // Suggested folder in which the file should be saved.
}

// SaveFiles opens a filechooser for selecting where to save one or more files.
// The chooser will use the supplied title as it's name.
func SaveFiles(parentWindow, title string, options *SaveFilesOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{
		"modal": dbus.MakeVariant(!options.NotModal),
	}

	if options.AcceptLabel != "" {
		data["accept_label"] = dbus.MakeVariant(options.AcceptLabel)
	}

	if options.CurrentFolder != "" {
		nullTerminatedByteString := []byte(options.CurrentFolder + "\000")
		data["current_folder"] = dbus.MakeVariant(nullTerminatedByteString)
	}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(fileChooserCallName+".SaveFiles", 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}
