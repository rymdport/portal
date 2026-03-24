package usb

import (
	"github.com/godbus/dbus/v5"

	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
	"github.com/rymdport/portal/internal/request"
)

const acquireDevicesCallName = interfaceName + ".AcquireDevices"

// AcquireDeviceOptions contains a device ID and access options for AcquireDevices.
// Each element of the devices array contains the device ID, and an access option vardict with the following keys.
type AcquireDeviceOptions struct {
	ID       string // Device ID.
	Writable bool   // Whether the device will be opened in read-write or read-only mode. Default: False.
}

type dbusAcquireDevice struct {
	ID      string
	Options map[string]dbus.Variant
}

// AcquireDevices requests to acquire (i.e. open) the given device nodes.
// The process of acquiring is finished by calling FinishAcquireDevices after the request emitted a Success response.
//
// The org.freedesktop.portal.Request::Response signal is emitted without any extra information.
func AcquireDevices(parentWindow string, options []AcquireDeviceOptions) (dbus.ObjectPath, error) {
	devices := make([]dbusAcquireDevice, len(options))
	for i, dev := range options {
		opts := map[string]dbus.Variant{}
		if dev.Writable {
			opts["writable"] = convert.FromBool(true)
		}

		devices[i] = dbusAcquireDevice{
			ID:      dev.ID,
			Options: opts,
		}
	}

	data := map[string]dbus.Variant{}

	result, err := apis.Call(acquireDevicesCallName, parentWindow, devices, data)
	if err != nil {
		return "", err
	}

	status, _, err := request.OnSignalResponse(result.(dbus.ObjectPath))
	if err != nil {
		return "", err
	} else if status >= request.Cancelled {
		return "", nil
	}

	return result.(dbus.ObjectPath), nil
}
