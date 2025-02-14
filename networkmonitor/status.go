package networkmonitor

import (
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
	available, err := apis.Call(getAvailableCallName)
	if err != nil {
		return false, err
	}

	return available.(bool), err
}

// GetMetered returns whether the network is considered metered.
// That is, whether the system as traffic flowing through the default connection
// that is subject ot limitations by service providers.
func GetMetered() (bool, error) {
	metered, err := apis.Call(getMeteredCallName)
	if err != nil {
		return false, err
	}

	return metered.(bool), err
}

// GetConnectivity returns more detailed information about the host’s network connectivity.
func GetConnectivity() (Connectivity, error) {
	connection, err := apis.Call(getConnectivityCallName)
	if err != nil {
		return 0, err
	}

	return connection.(Connectivity), err
}

// GetStatus returns availabilty, metered and connectivity status all at once.
func GetStatus() (*Status, error) {
	result, err := apis.Call(getStatusCallName)
	if err != nil {
		return nil, err
	}

	status := result.(map[string]any)
	return &Status{
		Available:  status["available"].(bool),
		Metered:    status["metered"].(bool),
		Connection: status["connectivity"].(Connectivity),
	}, nil
}

// CanReach returns whether the given hostname is believed to be reachable.
func CanReach(hostname string, port uint32) (bool, error) {
	reachable, err := apis.Call(canReachCallname, hostname, port)
	if err != nil {
		return false, err
	}

	return reachable.(bool), err
}
