package convert

import (
	"reflect"

	"github.com/godbus/dbus/v5"
)

var boolSignature = dbus.SignatureOfType(reflect.TypeOf(false))

// FromBool is a fast converter from a boolean to its dbus representation.
func FromBool(input bool) dbus.Variant {
	return dbus.MakeVariantWithSignature(input, boolSignature)
}
