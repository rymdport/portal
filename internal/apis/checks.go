package apis

import (
	"errors"

	"github.com/godbus/dbus/v5"
)

func checkDbusCompatibilityWitArgs(conn *dbus.Conn, args ...any) error {
	// Sending file through DBus must be tested before proceeding
	// See https://pkg.go.dev/github.com/godbus/dbus?utm_source=godoc#hdr-Unix_FD_passing
	if anyOf[any](args, isFileDescriptor) && !conn.SupportsUnixFDs() {
		return errors.New("DBus connection does not support passing unix file descriptors")
	}
	return nil
}

func isFileDescriptor(v any) bool {
	_, ok := v.(dbus.UnixFD)
	return ok
}

func anyOf[T any](values []T, assertion func(T) bool) bool {
	for _, v := range values {
		if assertion(v) {
			return true
		}
	}
	return false
}
