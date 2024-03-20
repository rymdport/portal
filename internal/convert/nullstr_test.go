package convert

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestToNullTerminated(t *testing.T) {
	input := "test"

	actual := ToNullTerminatedValue(input)
	expect := dbus.MakeVariant([]byte{'t', 'e', 's', 't', '\000'})
	if actual.Signature() != expect.Signature() {
		t.Fatalf("Expected %v, got %v", expect, actual)
	}
}

func BenchmarkToNullTerminated(b *testing.B) {
	var variant dbus.Variant

	for i := 0; i < b.N; i++ {
		variant = ToNullTerminatedValue("long_input_string")
	}

	benchVariant = variant
}
