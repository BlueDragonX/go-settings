package settings

import (
	"testing"
)

func TestParseSize(t *testing.T) {
	type test struct {
		str  string
		size int64
	}

	tests := []test{
		// conversions
		{"32", 32},
		{"56b", 56},
		{"14k", 14336},
		{"63kb", 64512},
		{"19m", 19922944},
		{"128mb", 134217728},
		{"15g", 16106127360},
		{"46gb", 49392123904},
		{"2t", 2199023255552},
		{"3t", 3298534883328},

		// case and whitespace
		{"14B", 14},
		{"92KB", 94208},
		{"92kB", 94208},
		{"92Kb", 94208},
		{" 92Kb", 94208},
		{" 92 Kb", 94208},
		{" 92  Kb ", 94208},
		{" 92 Kb", 94208},
		{"92 Kb ", 94208},
	}

	for _, test := range tests {
		if b, err := ParseSize(test.str); err != nil {
			t.Errorf("error parsing %s: %s", test.str, err)
		} else if b != test.size {
			t.Errorf("size is invalid: %d != %d", b, test.size)
		}
	}

	errors := []string{
		"15.0",
		"15 bits",
		"not a number",
	}

	for _, str := range errors {
		if _, err := ParseSize(str); err == nil {
			t.Errorf("no error parsing %s", str)
		}
	}
}
