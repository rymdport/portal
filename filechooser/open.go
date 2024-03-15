package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

const openFileCallName = fileChooserCallName + ".OpenFile"

// OpenFileOptions contains the options for how files are to be selected.
type OpenFileOptions struct {
	AcceptLabel   string // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool   // Whether the dialog should not be modal.
	Multiple      bool   // Whether multiple files can be selected or not.
	Directory     bool   // Whether to select for folders instead of files.
	CurrentFolder string // Suggested folder from which the files should be opened.
}

// OpenFile opens a filechooser for selecting a file to open.
// The chooser will use the supplied title as it's name.
func OpenFile(parentWindow, title string, options *OpenFileOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{
		"modal":     dbus.MakeVariant(!options.NotModal),
		"multiple":  dbus.MakeVariant(options.Multiple),
		"directory": dbus.MakeVariant(options.Directory),
	}

	if options.AcceptLabel != "" {
		data["accept_label"] = dbus.MakeVariant(options.AcceptLabel)
	}

	if options.CurrentFolder != "" {
		data["current_folder"] = convert.ToNullTerminatedString(options.CurrentFolder)
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(openFileCallName, 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}
