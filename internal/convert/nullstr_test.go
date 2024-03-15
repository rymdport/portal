package convert

import (
	"bytes"
	"testing"
)

func TestToNullTerminated(t *testing.T) {
	input := "test"

	got := ToNullTerminated(input)
	expect := []byte{'t', 'e', 's', 't', '\000'}
	if !bytes.Equal(got, expect) {
		t.Fatalf("Got %v, expected %v", got, expect)
	}
}

var benchResult []byte

func BenchmarkToNullTerminated(b *testing.B) {
	var result []byte

	for i := 0; i < b.N; i++ {
		result = ToNullTerminated("long_input_string")
	}

	benchResult = result
}
