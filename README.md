<a href="https://pkg.go.dev/fyne.io/fyne/v2?tab=doc" title="Go API Reference" rel="nofollow"><img src="https://img.shields.io/badge/go-documentation-blue.svg?style=flat" alt="Go API Reference"></a>
[![Analysis](https://github.com/rymdport/portal/actions/workflows/analysis.yml/badge.svg)](https://github.com/rymdport/portal/actions/workflows/analysis.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rymdport/portal)](https://goreportcard.com/report/github.com/rymdport/portal)

# Portal

Portal is a Go wrapper around the [XDG Desktop Protocol](https://flatpak.github.io/xdg-desktop-portal/) over DBUS.
This allows the code to request information from the user through the help of the desktop environment even when running inside a sandbox like Flatpak.  
As such, it is possible to easily open file dialogs, links and send notifications to the user in a way that integrates well with the desktop environment.

When running inside a sandbox, this runtime request of permissions allows the code to do various things without having to grant more permissions to the sandbox.
The portal APIs should also work good even when used outside of a sandboxed environment.

The goal of this project is to be a toolkit agnostic package for Go graphical user interface toolkits, and any other user case that needs it, to use the portals for improved Flatpak support etc.

## Supported Portal APIs

The following APIs are partially or completely implemented:

- [x] [Notification](https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Notification.html)
- [x] [OpenURI](https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.OpenURI.html)
- [x] [FileChooser](https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.FileChooser.html)

## Integrations into other projects

This section is meant as a reference to where this project is being used. Feel free to add yours if desired.

- As of the [v2.5.0](https://github.com/fyne-io/fyne/releases/tag/v2.5.0) release of [Fyne](https://fyne.io), this project is used when compiling with `-tags flatpak`. Parts of this projects might be used by default in the future. 

## Contributing

Contributions are strongly appreciated. Everything from creating bug reports to contributing code will help the project a lot, so please feel free to help in any way, shape, or form that you feel comfortable with.

## License
- Portal is licensed under ` Apache License Version 2.0` and will forever continue to be free and open source.
