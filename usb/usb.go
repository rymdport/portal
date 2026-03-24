// Package usb lets sandboxed applications monitor and request access to connected USB devices.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Usb.html.
package usb

import "github.com/rymdport/portal/internal/apis"

const interfaceName = apis.CallBaseName + ".Usb"
