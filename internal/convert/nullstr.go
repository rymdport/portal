package convert

import (
	"reflect"

	"github.com/godbus/dbus/v5"
)

var bytesSignature = dbus.SignatureOfType(reflect.TypeOf([]byte{}))

// ToNullTerminated connverts a regular string into a null terminated byte string.
func ToNullTerminated(input string) dbus.Variant {
	terminated := make([]byte, len(input)+1)
	copy(terminated, input)
	return dbus.MakeVariantWithSignature(terminated, bytesSignature)
}
