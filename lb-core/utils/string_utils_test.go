package utils

import "testing"

func TestSubstringLast(t *testing.T) {
	tests := []struct {
		name        string
		value, part string
		before      string
		after       string
	}{
		{name: "last occurrence", value: "a/b/c", part: "/", before: "a/b", after: "c"},
		{name: "missing occurrence", value: "filename", part: "/", before: "filename", after: ""},
		{name: "suffix", value: "filename/", part: "/", before: "filename", after: ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			before, after := SubstringLast(test.value, test.part)
			if before != test.before || after != test.after {
				t.Fatalf("SubstringLast(%q, %q) = (%q, %q), want (%q, %q)", test.value, test.part, before, after, test.before, test.after)
			}
		})
	}
}
