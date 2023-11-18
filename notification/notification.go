package notification

import (
	"strconv"

	"github.com/godbus/dbus/v5"
	"github.com/rymdport/portal"
)

const notificationCallName = portal.CallBaseName + ".Notification"

// Priority is the priroity of a notification.
type Priority = string

const (
	Low    Priority = "low"
	Normal Priority = "normal"
	High   Priority = "high"
	Urgent Priority = "urgent"
)

// Notification holds the content to send with the notification.
type Notification struct {
	Title    string
	Body     string
	Icon     string
	Priority Priority
}

// AddNotification sends a notification using org.freedesktop.portal.Notification.AddNotification.
func AddNotification(id uint, content *Notification) error {
	bus, err := dbus.SessionBus() // shared connection, don't close.
	if err != nil {
		return err
	}

	data := map[string]dbus.Variant{
		"title":    dbus.MakeVariant(content.Title),
		"body":     dbus.MakeVariant(content.Body),
		"icon":     dbus.MakeVariant(content.Icon),
		"priority": dbus.MakeVariant(content.Priority),
	}

	obj := bus.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(notificationCallName+".AddNotification", 0, strconv.FormatUint(uint64(id), 10), data)
	return call.Err
}

// RemoveNotification removes the notification with the corresponding id.
func RemoveNotification(id uint) error {
	bus, err := dbus.SessionBus() // shared connection, don't close.
	if err != nil {
		return err
	}

	obj := bus.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(notificationCallName+".RemoveNotification", 0, strconv.FormatUint(uint64(id), 10))
	return call.Err
}