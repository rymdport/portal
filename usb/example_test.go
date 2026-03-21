package usb_test

import (
	"encoding/json"
	"os"

	"github.com/rymdport/portal/usb"
)

func ExampleEnumerateDevices() {
	devices, err := usb.EnumerateDevices()
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	for _, dev := range devices {
		props := make(map[string]any, len(dev.Properties))
		for k, v := range dev.Properties {
			props[k] = v.Value()
		}

		if err := enc.Encode(map[string]any{
			"ID":         dev.ID,
			"Parent":     dev.Parent,
			"Readable":   dev.Readable,
			"Writable":   dev.Writable,
			"DeviceFile": dev.DeviceFile,
			"Properties": props,
		}); err != nil {
			panic(err)
		}
	}
}

func ExampleAcquireDevices() {
	handle, err := usb.AcquireDevices("", []usb.AcquireDeviceOptions{{ID: "device-id", Writable: false}})
	if err != nil {
		panic(err)
	}

	result, err := usb.FinishAcquireDevices(handle)
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
	for _, r := range result.Results {
		fd := ^uintptr(0)
		if r.File != nil {
			fd = r.File.Fd()
		}

		if err := enc.Encode(map[string]any{
			"ID":       r.ID,
			"Success":  r.Success,
			"Fd":       fd,
			"Error":    r.Error,
			"Finished": result.Finished,
		}); err != nil {
			panic(err)
		}
	}
}

func ExampleReleaseDevices() {
	if err := usb.ReleaseDevices([]string{"device-id"}); err != nil {
		panic(err)
	}
}

func ExampleCreateSession() {
	sess, err := usb.CreateSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	done := make(chan struct{})
	sess.SetOnClosed(func(err error) {
		if err != nil {
			panic(err)
		}
		close(done)
	})

	enc := json.NewEncoder(os.Stdout)
	if err := sess.SetOnDeviceEvents(func(events []usb.DeviceEvent) {
		for _, ev := range events {
			device := make(map[string]any, len(ev.Device))
			for k, v := range ev.Device {
				device[k] = v.Value()
			}

			if err := enc.Encode(map[string]any{
				"Action": ev.Action,
				"ID":     ev.ID,
				"Device": device,
			}); err != nil {
				panic(err)
			}
		}
	}); err != nil {
		panic(err)
	}

	<-done
}
