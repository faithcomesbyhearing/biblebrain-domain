package storage

import (
	"testing"
)

func TestNew(t *testing.T) {
	type test struct {
		base     string
		bible    string
		fileset  string
		filename string
		want     string
	}

	tests := []test{
		{want: ""},
		{base: "a", want: "a"},
		{base: "a", bible: "b", want: "a/b"},
		{base: "a", bible: "b", fileset: "c", want: "a/b/c"},
		{base: "a", bible: "b", fileset: "c", filename: "d", want: "a/b/c/d"},
	}

	for _, tc := range tests {
		got := NewPrefix(tc.base, tc.bible, tc.fileset, tc.filename)
		if tc.want != got.String() {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestParse(t *testing.T) {
	type test struct {
		str  string
		want Prefix
	}

	tests := []test{
		{str: "", want: NewPrefix("", "", "", "")},
		{str: "a", want: NewPrefix("a", "", "", "")},
		{str: "a/b", want: NewPrefix("a", "b", "", "")},
		{str: "a/b/c", want: NewPrefix("a", "b", "c", "")},
		{str: "a/b/c/d", want: NewPrefix("a", "b", "c", "d")},
		{str: "a/b/c/d/e", want: NewPrefix("a", "b", "c", "d")},
	}

	for _, tc := range tests {
		got := Parse(tc.str)
		if tc.want != got {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestString(t *testing.T) {
	type test struct {
		P    Prefix
		want string
	}

	tests := []test{
		{P: NewPrefix("", "", "", ""), want: ""},
		{P: NewPrefix("a", "", "", ""), want: "a"},
		{P: NewPrefix("a", "b", "", ""), want: "a/b"},
		{P: NewPrefix("a", "b", "c", "d"), want: "a/b/c/d"},
	}

	for _, tc := range tests {
		got := tc.P.String()
		if tc.want != got {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func TestCompare(t *testing.T) {
	type test struct {
		P1   Prefix
		P2   Prefix
		want int
	}

	tests := []test{
		{P1: NewPrefix("", "", "", ""), P2: NewPrefix("", "", "", ""), want: 0},
		{P1: NewPrefix("a", "", "", ""), P2: NewPrefix("a", "", "", ""), want: 0},
		{P1: NewPrefix("a", "", "", ""), P2: NewPrefix("b", "", "", ""), want: -1},
		{P1: NewPrefix("b", "", "", ""), P2: NewPrefix("a", "", "", ""), want: +1},
		{P1: NewPrefix("a", "a", "", ""), P2: NewPrefix("a", "a", "", ""), want: 0},
		{P1: NewPrefix("a", "a", "", ""), P2: NewPrefix("a", "b", "", ""), want: -1},
		{P1: NewPrefix("a", "b", "", ""), P2: NewPrefix("a", "a", "", ""), want: +1},
		{P1: NewPrefix("a", "b", "c", ""), P2: NewPrefix("a", "b", "c", ""), want: 0},
		{P1: NewPrefix("a", "b", "c", ""), P2: NewPrefix("a", "b", "d", ""), want: -1},
		{P1: NewPrefix("a", "b", "d", ""), P2: NewPrefix("a", "b", "c", ""), want: +1},
		{P1: NewPrefix("a", "b", "c", "d"), P2: NewPrefix("a", "b", "c", "d"), want: 0},
		{P1: NewPrefix("a", "b", "c", "d"), P2: NewPrefix("a", "b", "c", "e"), want: -1},
		{P1: NewPrefix("a", "b", "c", "e"), P2: NewPrefix("a", "b", "c", "d"), want: +1},
	}

	for _, tc := range tests {
		got := Compare(tc.P1, tc.P2)
		if tc.want != got {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}
