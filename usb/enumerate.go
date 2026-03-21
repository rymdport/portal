package usb

import (
	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
)

const enumerateDevicesCallName = interfaceName + ".EnumerateDevices"

// DeviceProperties holds a list of udev properties that a device has.
// These properties are not parsed in any way by the portal, it is up to apps to parse them.
type DeviceProperties = map[string]dbus.Variant

// Device represents a connected USB device.
// Each element of the devices array contains the device ID, and a device information vardict.
type Device struct {
	ID         string           // Device ID.
	Parent     string           // Device ID of the parent device.
	Readable   bool             // Whether the device can be opened for reading with AcquireDevices. If not present, then it should be assumed to be false.
	Writable   bool             // Whether the device can be opened for writing with AcquireDevices. If not present, then it should be assumed to be false.
	DeviceFile string           // A string path to the device node inside the /dev filesystem.
	Properties DeviceProperties // A list of udev properties that this device has.
}

// EnumerateDevices enumerates all connected USB devices that this application has permission to see.
func EnumerateDevices() ([]Device, error) {
	options := map[string]dbus.Variant{}

	result, err := apis.Call(enumerateDevicesCallName, options)
	if err != nil {
		return nil, err
	}

	rawDevices := result.([][]any)
	devices := make([]Device, len(rawDevices))
	for i, raw := range rawDevices {
		id := raw[0].(string)
		info := raw[1].(map[string]dbus.Variant)

		dev := Device{ID: id}
		if v, ok := info["parent"]; ok {
			dev.Parent = v.Value().(string)
		}
		if v, ok := info["readable"]; ok {
			dev.Readable = v.Value().(bool)
		}
		if v, ok := info["writable"]; ok {
			dev.Writable = v.Value().(bool)
		}
		if v, ok := info["device-file"]; ok {
			dev.DeviceFile = v.Value().(string)
		}
		if v, ok := info["properties"]; ok {
			dev.Properties = v.Value().(map[string]dbus.Variant)
		}

		devices[i] = dev
	}

	return devices, nil
}
