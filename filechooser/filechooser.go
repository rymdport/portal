// Package filechooser allows sandboxed applications to ask the user for access to files outside the sandbox. The portal backend will present the user with a file chooser dialog.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.FileChooser.html.
package filechooser

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/request"
)

const interfaceName = apis.CallBaseName + ".FileChooser"

func readURIFromResponse(path dbus.ObjectPath) ([]string, error) {
	status, results, err := request.OnSignalResponse(path)
	if err != nil {
		return nil, err
	} else if status == request.Cancelled {
		return nil, nil
	}

	uris := results["uris"].Value().([]string)
	return uris, nil
}
