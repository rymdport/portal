[![Go API Reference](https://img.shields.io/badge/go-documentation-blue.svg?style=flat)](https://pkg.go.dev/github.com/rymdport/portal)
[![Analysis](https://github.com/rymdport/portal/actions/workflows/analysis.yml/badge.svg)](https://github.com/rymdport/portal/actions/workflows/analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rymdport/portal)](https://goreportcard.com/report/github.com/rymdport/portal)

# Portal

Portal is a Go wrapper around the [XDG Desktop Protocol](https://flatpak.github.io/xdg-desktop-portal/) over DBUS.
This allows the code to request information from the user through the help of the desktop environment even when running inside a sandbox like Flatpak.  
As such, it is possible to easily open file dialogs, links and send notifications to the user in a way that integrates well with the desktop environment.

When running inside a sandbox, this runtime request of permissions allows the code to do various things without having to grant more permissions to the sandbox.
The portal APIs should also work good even when used outside of a sandboxed environment.

The goal of this project is to be a toolkit agnostic package for Go graphical user interface toolkits, and any other user case that needs it, to use the portals for improved Flatpak support etc.

## Example

The following example showcases how a file chooser can be opened for selecting one or more files.


```go
package main

import (
	"fmt"
	"log"

	"github.com/rymdport/portal/filechooser"
)

func main() {
	options := filechooser.OpenFileOptions{Multiple: true}
	files, err := filechooser.OpenFile("", "Select files", &options)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(files)
}
```

## Supported Portal APIs

The lsit below contains all of the portal APIs available as of 2024-03-14. Checked boxes are partially or completely implemented within this project.

- [ ] Account
- [ ] Background
- [ ] Camera
- [ ] Clipboard
- [ ] Device
- [ ] Documents
- [ ] Dynamic Launcher
- [ ] Email
- [x] File Chooser
- [ ] File Transfer
- [ ] Game Mode
- [ ] Global Shortcuts
- [ ] Inhibit
- [ ] Input Capture
- [ ] Location
- [ ] Memory Monitor
- [ ] Network Monitor
- [x] Notification
- [x] OpenURI
- [ ] Power Profile Monitor
- [ ] Print
- [ ] Proxy Resolver
- [ ] Realtime
- [ ] Remote Desktop
- [ ] Request
- [ ] ScreenCast
- [ ] Screenshot
- [ ] Secret
- [ ] Session
- [x] Settings
- [ ] Trash
- [ ] Wallpaper


## Used by other projects

This section is meant as a reference to where this project is being used. Feel free to add yours if desired.

- As of the [v2.5.0](https://github.com/fyne-io/fyne/releases/tag/v2.5.0) release of [Fyne](https://fyne.io), this project is used when compiling with `-tags flatpak`. Parts of this projects might be used by default in the future. 

## Contributing

Contributions are strongly appreciated. Everything from creating bug reports to contributing code will help the project a lot, so please feel free to help in any way, shape, or form that you feel comfortable with.

## License
- Portal is licensed under ` Apache License Version 2.0` and will forever continue to be free and open source.
