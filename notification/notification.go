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

// Content holds the content to send with the notification.
type Content struct {
	Title    string
	Body     string
	Icon     string
	Priority Priority
}

// Add sends a notification using org.freedesktop.portal.Notification.Add.
func Add(id uint, content *Content) error {
	bus, err := dbus.SessionBus() // shared connection, don't close.
	if err != nil {
		return err
	}

	data := map[string]dbus.Variant{
		"title": dbus.MakeVariant(content.Title),
		"body":  dbus.MakeVariant(content.Body),
		"icon":  dbus.MakeVariant(content.Icon),
	}

	// Only add the priority field when it is set.
	if content.Priority != "" {
		data["priority"] = dbus.MakeVariant(content.Priority)
	}

	obj := bus.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(notificationCallName+".AddNotification", 0, strconv.FormatUint(uint64(id), 10), data)
	return call.Err
}

// Remove removes the notification with the corresponding id.
func Remove(id uint) error {
	bus, err := dbus.SessionBus() // shared connection, don't close.
	if err != nil {
		return err
	}

	obj := bus.Object(portal.ObjectName, portal.ObjectPath)
	call := obj.Call(notificationCallName+".RemoveNotification", 0, strconv.FormatUint(uint64(id), 10))
	return call.Err
}
