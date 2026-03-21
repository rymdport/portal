package usb

import (
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal/internal/apis"
	"github.com/rymdport/portal/internal/convert"
)

const finishAcquireDevicesCallName = interfaceName + ".FinishAcquireDevices"

// AcquireResult contains the result of acquiring a single device.
// Each element of the results array contains the device ID, and a vardict with the result.
type AcquireResult struct {
	ID      string   // Device ID.
	Success bool     // Whether the device access was successful or not.
	File    *os.File // The file descriptor representing the device. The file descriptor is meant to be passed on to the USB library. Only present if this was a successful device access.
	Error   string   // Error message describing why accessing the device was not successful. Only present if this was an failed device access.
}

// FinishAcquireDevicesResult contains the results of FinishAcquireDevices.
type FinishAcquireDevicesResult struct {
	Results  []AcquireResult // Array of device ids, and the result of the access.
	Finished bool            // Whether all device results were reported.
}

// FinishAcquireDevices retrieves the file descriptors of the devices requested during AcquireDevices.
// The file descriptors remain usable until the underlying device is removed, they are released with ReleaseDevices,
// the D-Bus connections is closed, or the portal revokes the file descriptor which can happen at any point.
// Devices which are not needed anymore should be passed to ReleaseDevices.
//
// If not all devices could be send, Finished will be false and FinishAcquireDevices must be called again
// until Finished is true, before calling AcquireDevices again.
//
// This method can only be called once for a given token, and only after calling AcquireDevices.
func FinishAcquireDevices(handle dbus.ObjectPath) (*FinishAcquireDevicesResult, error) {
	data := map[string]dbus.Variant{}

	var (
		rawResults [][]any
		finished   bool
	)
	if err := apis.CallStore(finishAcquireDevicesCallName, []any{&rawResults, &finished}, handle, data); err != nil {
		return nil, err
	}

	results := make([]AcquireResult, len(rawResults))
	for i, raw := range rawResults {
		id := raw[0].(string)
		info := raw[1].(map[string]dbus.Variant)

		r := AcquireResult{ID: id}
		if v, ok := info["success"]; ok {
			r.Success = v.Value().(bool)
		}
		if v, ok := info["fd"]; ok {
			fd, err := convert.UnixFDToUintptr(v.Value().(dbus.UnixFD))
			if err == nil {
				r.File = os.NewFile(fd, id)
			}
		}
		if v, ok := info["error"]; ok {
			r.Error = v.Value().(string)
		}

		results[i] = r
	}

	return &FinishAcquireDevicesResult{
		Results:  results,
		Finished: finished,
	}, nil
}
