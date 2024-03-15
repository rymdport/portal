package convert

import "github.com/godbus/dbus/v5"

// ToNullTerminatedString connverts a regular string into a null terminated dbus variant string.
func ToNullTerminatedString(input string) dbus.Variant {
	return dbus.MakeVariant([]byte(input + "\000"))
}
