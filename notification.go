package portal

import (
	"strconv"

	"github.com/godbus/dbus/v5"
)

// Notification holds the content to send with the notification.
type Notification struct {
	Title string
	Body  string
	Icon  string
}

// AddNotification sends a notification using org.freedesktop.portal.Notification.AddNotification.
func AddNotification(id uint, content *Notification) error {
	bus, err := dbus.SessionBus() // shared connection, don't close.
	if err != nil {
		return err
	}

	data := map[string]dbus.Variant{
		"title": dbus.MakeVariant(content.Title),
		"body":  dbus.MakeVariant(content.Body),
		"icon":  dbus.MakeVariant(content.Icon),
	}

	obj := bus.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	call := obj.Call("org.freedesktop.portal.Notification.AddNotification", 0, strconv.FormatUint(uint64(id), 10), data)
	return call.Err
}

// RemoveNotification removes the notification with the corresponding id.
func RemoveNotification(id uint) error {
	bus, err := dbus.SessionBus() // shared connection, don't close.
	if err != nil {
		return err
	}

	obj := bus.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	call := obj.Call("org.freedesktop.portal.Notification.RemoveNotification", 0, strconv.FormatUint(uint64(id), 10))
	return call.Err
}
