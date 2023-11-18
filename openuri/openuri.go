package openuri

import (
	"github.com/fredbi/uri"
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const openURICallName = portal.CallBaseName + ".OpenURI"

// OpenURI opens the given URI in the corresponding application.
func OpenURI(uri uri.URI) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	parentWindow := ""
	data := map[string]dbus.Variant{}

	obj := conn.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(openURICallName+".OpenURI", 0, parentWindow, uri.String(), data)
	return call.Err
}
