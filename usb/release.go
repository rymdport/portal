package usb

import (
	"github.com/godbus/dbus/v5"

	"github.com/rymdport/portal/internal/apis"
)

const releaseDevicesCallName = interfaceName + ".ReleaseDevices"

// ReleaseDevices releases previously acquired devices.
// The file descriptors of those devices might become unusable as a result of this.
//
// Each element of the devices array contains the device ID of the device.
// There are no supported keys in the options vardict.
func ReleaseDevices(deviceIDs []string) error {
	data := map[string]dbus.Variant{}

	return apis.CallWithoutResult(releaseDevicesCallName, deviceIDs, data)
}
