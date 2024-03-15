// Package openuri allows sandboxed applications to open URIs (e.g. a http: link to the applications homepage) under the control of the user.
// Upstream API documentation can be found at https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.OpenURI.html.
package openuri

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const openURICallName = portal.CallBaseName + ".OpenURI"

// OpenURI opens the given URI in the corresponding application.
func OpenURI(parentWindow, uri string) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	data := map[string]dbus.Variant{}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(openURICallName+".OpenURI", 0, parentWindow, uri, data)
	return call.Err
}
