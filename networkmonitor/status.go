package networkmonitor

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const (
	getAvailableCallName    = interfaceName + ".GetAvailable"
	getMeteredCallName      = interfaceName + ".GetMetered"
	getConnectivityCallName = interfaceName + ".GetConnectivity"
	getStatusCallName       = interfaceName + ".GetStatus"
	canReachCallname        = interfaceName + ".CanReach"
)

// Connectivity specifies detailed information about network connectivity.
type Connectivity = uint32

const (
	LocalOnly     Connectivity = 1 // The host is not configured with a route to the internet.
	Limited       Connectivity = 2 // The host is connected to a network, but can't reach the full internet.
	CaptivePortal Connectivity = 3 // The host is behind a captive portal and cannot reach the full internet.
	FullNetwork   Connectivity = 4 // The host connected to a network, and can reach the full internet.
)

// Status contains common network status values.
type Status struct {
	Available  bool
	Metered    bool
	Connection Connectivity
}

// GetAvailable returns whether the network is considered available.
// That is, whether the system as a default route for at least one of IPv4 or IPv6.
func GetAvailable() (bool, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return false, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(getAvailableCallName, 0)
	if call.Err != nil {
		return false, call.Err
	}

	var available bool
	err = call.Store(&available)
	return available, err
}

// GetMetered returns whether the network is considered metered.
// That is, whether the system as traffic flowing through the default connection
// that is subject ot limitations by service providers.
func GetMetered() (bool, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return false, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(getMeteredCallName, 0)
	if call.Err != nil {
		return false, call.Err
	}

	var metered bool
	err = call.Store(&metered)
	return metered, err
}

// GetConnectivity returns more detailed information about the host’s network connectivity.
func GetConnectivity() (Connectivity, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return 0, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(getConnectivityCallName, 0)
	if call.Err != nil {
		return 0, call.Err
	}

	var connectivity Connectivity
	err = call.Store(&connectivity)
	return connectivity, err
}

// GetStatus returns availabilty, metered and connectivity status all at once.
func GetStatus() (*Status, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return nil, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(getStatusCallName, 0)
	if call.Err != nil {
		return nil, call.Err
	}

	var status map[string]any
	err = call.Store(&status)
	if err != nil {
		return nil, err
	}

	return &Status{
		Available:  status["available"].(bool),
		Metered:    status["metered"].(bool),
		Connection: status["connectivity"].(Connectivity),
	}, nil
}

// CanReach returns whether the given hostname is believed to be reachable.
func CanReach(hostname string, port uint32) (bool, error) {
	conn, err := dbus.SessionBus() // Shared connection, don't close.
	if err != nil {
		return false, err
	}

	obj := conn.Object(apis.ObjectName, apis.ObjectPath)
	call := obj.Call(canReachCallname, 0, hostname, port)
	if call.Err != nil {
		return false, call.Err
	}

	var reachable bool
	err = call.Store(&reachable)
	return reachable, err
}
