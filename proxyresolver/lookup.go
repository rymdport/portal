package proxyresolver

import "github.com/rymdport/portal/internal/apis"

const lookupCallName = interfaceName + ".Lookup"

// Lookup looks up which proxy to use to connect to uri.
// The returned proxy uri are of the form protocol://[user[:password] AT host:port.
// The protocol can be http, rtsp, socks or another proxying protocol.
// direct:// is used when no proxy is needed.
func Lookup(uri string) ([]string, error) {
	proxies, err := apis.Call(lookupCallName, uri)
	if err != nil {
		return nil, err
	}

	return proxies.([]string), err
}
