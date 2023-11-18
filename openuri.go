package portal

import (
	"github.com/fredbi/uri"
	"github.com/godbus/dbus/v5"
)

const openURICallName = callBaseName + "OpenURI"

// OpenURI opens the given URI in the corresponding application.
func OpenURI(uri uri.URI) error {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return err
	}

	parentWindow := ""
	options := map[string]dbus.Variant{}

	obj := conn.Object(objectName, objectPath)
	call := obj.Call(openURICallName+".OpenURI", 0, parentWindow, uri.String(), options)
	return call.Err
}
