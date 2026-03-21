package main

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/rymdport/portal/usb"
)

func main() {
	deviceID := flag.String("device", "", "device ID")
	writable := flag.Bool("writable", false, "acquire device in read-write mode")

	flag.Parse()

	enc := json.NewEncoder(os.Stdout)

	switch flag.Arg(0) {
	case "enumerate":
		enumerate(enc)

	case "acquire":
		acquire(enc, *deviceID, *writable)

	case "release":
		release(*deviceID)

	case "monitor":
		monitor(enc)

	default:
		panic("unknown command: " + flag.Arg(0) + " (expected: enumerate, acquire, release, monitor)")
	}
}

func enumerate(enc *json.Encoder) {
	devices, err := usb.EnumerateDevices()
	if err != nil {
		panic(err)
	}

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

func acquire(enc *json.Encoder, deviceID string, writable bool) {
	handle, err := usb.AcquireDevices("", []usb.AcquireDeviceOptions{{ID: deviceID, Writable: writable}})
	if err != nil {
		panic(err)
	}

	result, err := usb.FinishAcquireDevices(handle)
	if err != nil {
		panic(err)
	}

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

func release(deviceID string) {
	if err := usb.ReleaseDevices([]string{deviceID}); err != nil {
		panic(err)
	}
}

func monitor(enc *json.Encoder) {
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
