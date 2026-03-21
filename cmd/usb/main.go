package main

import (
	"encoding/json"
	"maps"
	"os"

	"github.com/rymdport/portal/usb"
)

func main() {
	devices, err := usb.EnumerateDevices()
	if err != nil {
		panic(err)
	}

	enc := json.NewEncoder(os.Stdout)
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
}
