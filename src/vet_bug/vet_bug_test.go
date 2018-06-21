package vet_bug

import (
	"testing"
)

func TestAdd(t *testing.T) {
	inouts := []struct {
		a int
		b int
		r int
	}{
		{a: 0, b: 0, r: 0},
		{a: 1, b: 1, r: 2},
		{a: 3, b: 5, r: 8},
	}

	for _, inout := range inouts {
		AssertEqual(t, inout.r, Add(inout.a, inout.b))
	}
}
