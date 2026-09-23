package filechooser

import (
	"context"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/convert"
)

const (
	saveFileCallName  = interfaceName + ".SaveFile"
	saveFilesCallName = interfaceName + ".SaveFiles"
)

// SaveFileOptions contains the options for how a file is saved.
type SaveFileOptions struct {
	HandleToken   string      // A string that will be used as the last element of the handle. Must be a valid object path element.
	AcceptLabel   string      // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool        // Whether the dialog should not be modal.
	Filters       []*Filter   // Each item specifies a single filter to offer to the user.
	CurrentFilter *Filter     // Request that this filter be set by default at dialog creation.
	Choices       []*ComboBox // List of serialized combo boxes to add to the file chooser.
	CurrentFolder string      // Suggested folder in which the file should be saved.
	CurrentName   string      // Suggested name of the file.
}

// SaveFile opens a filechooser for selecting where to save a file.
// The chooser will use the supplied title as it's name.
func SaveFile(parentWindow, title string, options *SaveFileOptions) ([]string, error) {
	return SaveFileContext(context.Background(), parentWindow, title, options)
}

// SaveFileContext is SaveFile with a context.
func SaveFileContext(ctx context.Context, parentWindow, title string, options *SaveFileOptions) ([]string, error) {
	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}

	return callForURIs(ctx, userToken, saveFileCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		if options != nil {
			data["modal"] = convert.FromBool(!options.NotModal)

			if options.AcceptLabel != "" {
				data["accept_label"] = convert.FromString(options.AcceptLabel)
			}
			if len(options.Filters) > 0 {
				data["filters"] = dbus.MakeVariant(options.Filters)
			}
			if options.CurrentFilter != nil {
				data["current_filter"] = dbus.MakeVariant(options.CurrentFilter)
			}
			if len(options.Choices) > 0 {
				data["choices"] = dbus.MakeVariant(options.Choices)
			}
			if options.CurrentName != "" {
				data["current_name"] = convert.FromString(options.CurrentName)
			}
			if options.CurrentFolder != "" {
				data["current_folder"] = convert.ToNullTerminatedValue(options.CurrentFolder)
			}
		}

		return []any{parentWindow, title, data}
	})
}

// SaveFilesOptions contains the options for how files are saved.
type SaveFilesOptions struct {
	HandleToken   string      // A string that will be used as the last element of the handle. Must be a valid object path element.
	AcceptLabel   string      // Label for the accept button. Mnemonic underlines are allowed.
	NotModal      bool        // Whether the dialog should be modal.
	Choices       []*ComboBox // List of serialized combo boxes to add to the file chooser.
	CurrentFolder string      // Suggested folder in which the file should be saved.
	Files         []string    // An array of file names to be saved.
}

// SaveFiles opens a filechooser for selecting where to save one or more files.
// The chooser will use the supplied title as it's name.
func SaveFiles(parentWindow, title string, options *SaveFilesOptions) ([]string, error) {
	return SaveFilesContext(context.Background(), parentWindow, title, options)
}

// SaveFilesContext is SaveFiles with a context.
func SaveFilesContext(ctx context.Context, parentWindow, title string, options *SaveFilesOptions) ([]string, error) {
	userToken := ""
	if options != nil {
		userToken = options.HandleToken
	}

	return callForURIs(ctx, userToken, saveFilesCallName, func(token string) []any {
		data := map[string]dbus.Variant{
			"handle_token": convert.FromString(token),
		}
		if options != nil {
			data["modal"] = convert.FromBool(!options.NotModal)

			if options.AcceptLabel != "" {
				data["accept_label"] = convert.FromString(options.AcceptLabel)
			}
			if len(options.Choices) > 0 {
				data["choices"] = dbus.MakeVariant(options.Choices)
			}
			if options.CurrentFolder != "" {
				data["current_folder"] = convert.ToNullTerminatedValue(options.CurrentFolder)
			}
			if len(options.Files) > 0 {
				files := make([][]byte, len(options.Files))
				for i, file := range options.Files {
					files[i] = convert.FromStringToNullTerminated(file)
				}
				data["files"] = dbus.MakeVariant(files)
			}
		}

		return []any{parentWindow, title, data}
	})
}
