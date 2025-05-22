package apis

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func isNegative(v int) bool {
	return v < 0
}

func TestAnyOfIsNegativeTrue(t *testing.T) {
	testValues := []int{5, 7, -4, 6, 8}
	expected := true
	result := anyOf[int](testValues, isNegative)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestAnyOfIsNegativeFalse(t *testing.T) {
	testValues := []int{5, 7, 4, 6, 8}
	expected := false
	result := anyOf[int](testValues, isNegative)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestAnyOfIsNegativeEmpty(t *testing.T) {
	testValues := []int{}
	expected := false
	result := anyOf[int](testValues, isNegative)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestAnyOfIsNegativeFirst(t *testing.T) {
	testValues := []int{-5, 7, 4, 6, 8}
	expected := true
	result := anyOf[int](testValues, isNegative)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestAnyOfIsNegativeLast(t *testing.T) {
	testValues := []int{5, 7, 4, 6, -8}
	expected := true
	result := anyOf[int](testValues, isNegative)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestIsFDTrue(t *testing.T) {
	var tested dbus.UnixFD
	expected := true
	result := isFileDescriptor(tested)

	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestIsFDInt(t *testing.T) {
	var tested int
	expected := false
	result := isFileDescriptor(tested)

	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestIsFDInterface(t *testing.T) {
	var tested interface{}
	expected := false
	result := isFileDescriptor(tested)

	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestAnyOfIsFDTrue(t *testing.T) {
	var fd dbus.UnixFD
	var i interface{}

	testValues := []any{5, fd, i}
	expected := true
	result := anyOf[any](testValues, isFileDescriptor)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}

func TestAnyOfIsFDFalse(t *testing.T) {
	var i interface{}

	testValues := []any{5, -7, i}
	expected := false
	result := anyOf[any](testValues, isFileDescriptor)
	if expected != result {
		t.Errorf("Expecting %v, got %v", expected, result)
	}
}
