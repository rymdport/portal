package convert

import (
	"reflect"

	"github.com/godbus/dbus/v5"
)

var (
	boolSignature   = dbus.SignatureOfType(reflect.TypeOf(false))
	stringSignature = dbus.SignatureOfType(reflect.TypeOf(""))
	uint32Signature = dbus.SignatureOfType(reflect.TypeOf(uint32(0)))
)

// FromBool is a fast converter from a boolean to its dbus representation.
func FromBool(input bool) dbus.Variant {
	return dbus.MakeVariantWithSignature(input, boolSignature)
}

// FromBool is a fast converter from a string to its dbus representation.
func FromString(input string) dbus.Variant {
	return dbus.MakeVariantWithSignature(input, stringSignature)
}

// FromUint32 is a fast converter from a uint32 to its dbus representation.
func FromUint32(input uint32) dbus.Variant {
	return dbus.MakeVariantWithSignature(input, uint32Signature)
}
