package networkmonitor

import (
	"github.com/rymdport/portal/internal/apis"
)

// OnSignalChanged calls the passed function when the network configuration changes.
func OnSignalChanged(callback func()) error {
	signal, err := apis.ListenOnSignal(interfaceName, "changed")
	if err != nil {
		return err
	}

	for range signal {
		callback()
	}

	return nil
}
