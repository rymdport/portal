// Package usb provides access to USB devices through the portal.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Usb.html.
package usb

import "github.com/rymdport/portal/internal/apis"

const interfaceName = apis.CallBaseName + ".Usb"
