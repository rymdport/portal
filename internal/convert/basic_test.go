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

func TestFromString(t *testing.T) {
	actual := FromString("")
	expect := dbus.MakeVariant("")

	if actual != expect {
		t.Fatalf("Expected %v, got %v", expect, actual)
	}

	actual = FromString("testing")
	expect = dbus.MakeVariant("testing")

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

func BenchmarkFromString(b *testing.B) {
	variant := dbus.Variant{}

	input := "example"

	for i := 0; i < b.N; i++ {
		variant = FromString(input)
	}

	benchVariant = variant
}

func BenchmarkMakeVariantFromString(b *testing.B) {
	variant := dbus.Variant{}

	input := "example"

	for i := 0; i < b.N; i++ {
		variant = dbus.MakeVariant(input)
	}

	benchVariant = variant
}
