package convert

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestFromBool(t *testing.T) {
	actual := FromBool(false)
	expect := dbus.MakeVariant(false)

	if actual != expect {
		t.Fatalf("Expected %v, got %v", expect, actual)
	}

	actual = FromBool(true)
	expect = dbus.MakeVariant(true)

	if actual != expect {
		t.Fatalf("Expected %v, got %v", expect, actual)
	}
}

var benchVariant dbus.Variant

func BenchmarkFromBool(b *testing.B) {
	variant := dbus.Variant{}

	for i := 0; i < b.N; i++ {
		variant = FromBool(true)
	}

	benchVariant = variant
}

func BenchmarkMakeVariantFromBool(b *testing.B) {
	variant := dbus.Variant{}

	for i := 0; i < b.N; i++ {
		variant = dbus.MakeVariant(true)
	}

	benchVariant = variant
}
