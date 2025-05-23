package apis

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func isNegative(v int) bool {
	return v < 0
}

var testAnyOfWithIsNegative map[string]struct{
	input []int
	expected bool
} = map[string]struct{
	input []int
	expected bool
} {
	"isPresent":
		{
			input: []int{5, 7, -4, 6, 8},
			expected: true,
		},
	"isAbscent":
		{
			input: []int{5, 7, 4, 6, 8},
			expected: false,
		},
	"isEmpty":
		{
			input: []int{},
			expected: false,
		},
	"fisrt":{
			input: 
		[]int{-5, 7, 4, 6, 8},
		expected: true,},
	"last": {
		input: []int{5, 7, 4, 6, -8},
		expected: true,
	},
}

func TestAnyOfIsNegative(t *testing.T) {
	for name, test := range testAnyOfWithIsNegative {
		t.Run(name, func (t *testing.T) {
			t.Parallel()
			if got, expected := anyOf[int](test.input, isNegative), test.expected; got != expected {
				t.Errorf("anyOf[int](%q, isNegative) returned %t, expected %t", test.input, got, expected)
			}
		})
	}
}

func TestIsFD(t *testing.T) {
	var fd dbus.UnixFD
	var itfc  interface{}
	tests := map[string]struct{
		input any
		expected bool
	}{
		"unixDF": {
			input: fd,
			expected: true,
		},
		"int": {
			input: 5,
			expected: false,
		},
		"interface": {
			input: itfc,
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func (t *testing.T) {
			t.Parallel()
			if got, expected := isFileDescriptor(test.input), test.expected; got != expected {
				t.Errorf("isFileDescriptor(%v) returned %t, expected %t", test.input, got, expected)
			}
		})
	}
}

func TestAnyOfIsFDTrue(t *testing.T) {
	var fd dbus.UnixFD
	var i interface{}

	tests := map[string] struct{
		input []any
		expected bool
	}{
		"oneFd": {
			input: []any{5, fd, i},
			expected: true,
		},
		"noFd": {
			input: []any{5, -7, i},
			expected: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func (t *testing.T) {
			t.Parallel()
			if got, expected := anyOf[any](test.input, isFileDescriptor), test.expected; got != expected {
				t.Errorf("anyOf[any](%q, isFileDescriptor) returned %t, expected %t", test.input, got, expected)
			}
		})
	}
}