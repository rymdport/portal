package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

const (
	saveFileCallName  = fileChooserCallName + ".SaveFile"
	saveFilesCallName = fileChooserCallName + ".SaveFiles"
)

// SaveFileOptions contains the options for how a file is saved.
type SaveFileOptions struct {
	HandleToken   string   // A string that will be used as the last element of the handle. Must be a valid object path element.
	AcceptLabel   string   // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool     // Whether the dialog should not be modal.
	Filters       []Filter // Each item specifies a single filter to offer to the user.
	CurrentFilter *Filter  // Request that this filter be set by default at dialog creation.
	CurrentName   string   // Suggested name of the file.
	CurrentFolder string   // Suggested folder in which the file should be saved.
}

// SaveFile opens a filechooser for selecting where to save a file.
// The chooser will use the supplied title as it's name.
func SaveFile(parentWindow, title string, options *SaveFileOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{}

	if options != nil {
		data = map[string]dbus.Variant{
			"modal": convert.FromBool(!options.NotModal),
		}

		if options.HandleToken != "" {
			data["handle_token"] = convert.FromString(options.HandleToken)
		}

		if options.AcceptLabel != "" {
			data["accept_label"] = convert.FromString(options.AcceptLabel)
		}

		if len(options.Filters) > 0 {
			data["filters"] = dbus.MakeVariant(options.Filters)
		}

		if options.CurrentFilter != nil {
			data["current_filter"] = dbus.MakeVariant(options.CurrentFilter)
		}

		if options.CurrentName != "" {
			data["current_name"] = convert.FromString(options.CurrentName)
		}

		if options.CurrentFolder != "" {
			data["current_folder"] = convert.ToNullTerminated(options.CurrentFolder)
		}
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(saveFileCallName, 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}

// SaveFilesOptions contains the options for how files are saved.
type SaveFilesOptions struct {
	HandleToken   string   // A string that will be used as the last element of the handle. Must be a valid object path element.
	AcceptLabel   string   // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool     // Whether the dialog should be modal.
	Filters       []Filter // Each item specifies a single filter to offer to the user.
	CurrentFilter *Filter  // Request that this filter be set by default at dialog creation.
	CurrentFolder string   // Suggested folder in which the file should be saved.
}

// SaveFiles opens a filechooser for selecting where to save one or more files.
// The chooser will use the supplied title as it's name.
func SaveFiles(parentWindow, title string, options *SaveFilesOptions) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	data := map[string]dbus.Variant{}

	if options != nil {
		data = map[string]dbus.Variant{
			"modal": convert.FromBool(!options.NotModal),
		}

		if options.HandleToken != "" {
			data["handle_token"] = convert.FromString(options.HandleToken)
		}

		if options.AcceptLabel != "" {
			data["accept_label"] = convert.FromString(options.AcceptLabel)
		}

		if len(options.Filters) > 0 {
			data["filters"] = dbus.MakeVariant(options.Filters)
		}

		if options.CurrentFilter != nil {
			data["current_filter"] = dbus.MakeVariant(options.CurrentFilter)
		}

		if options.CurrentFolder != "" {
			data["current_folder"] = convert.ToNullTerminated(options.CurrentFolder)
		}
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(saveFilesCallName, 0, parentWindow, title, data)
	if call.Err != nil {
		return nil, call.Err
	}

	return readURIFromResponse(conn, call)
}
