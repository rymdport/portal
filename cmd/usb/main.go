package main

import (
	"encoding/json"
	"flag"
	"maps"
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
		devices, err := usb.EnumerateDevices()
		if err != nil {
			panic(err)
		}

		for _, dev := range devices {
			if err := enc.Encode(map[string]any{
				"ID":         dev.ID,
				"Parent":     dev.Parent,
				"Readable":   dev.Readable,
				"Writable":   dev.Writable,
				"DeviceFile": dev.DeviceFile,
				"Properties": maps.Collect(func(yield func(string, any) bool) {
					for k, v := range dev.Properties {
						if !yield(k, v.Value()) {
							return
						}
					}
				}),
			}); err != nil {
				panic(err)
			}
		}

	case "acquire":
		handle, err := usb.AcquireDevices("", []usb.AcquireDeviceOptions{{ID: *deviceID, Writable: *writable}})
		if err != nil {
			panic(err)
		}

		result, err := usb.FinishAcquireDevices(handle)
		if err != nil {
			panic(err)
		}

		for _, r := range result.Results {
			fd := -1
			if r.File != nil {
				fd = int(r.File.Fd())
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

	case "release":
		if err := usb.ReleaseDevices([]string{*deviceID}); err != nil {
			panic(err)
		}

	default:
		panic("unknown command: " + flag.Arg(0) + " (expected: enumerate, acquire, release)")
	}
}
