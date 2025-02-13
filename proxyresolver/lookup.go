package proxyresolver

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const lookupCallName = interfaceName + ".Lookup"

// Lookup looks up which proxy to use to connect to uri.
// The returned proxy uri are of the form protocol://[user[:password] AT host:port.
// The protocol can be http, rtsp, socks or another proxying protocol.
// direct:// is used when no proxy is needed.
func Lookup(uri string) ([]string, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(lookupCallName, 0, uri)
	if call.Err != nil {
		return nil, call.Err
	}

	var proxies []string
	err = call.Store(&proxies)
	return proxies, err
}
