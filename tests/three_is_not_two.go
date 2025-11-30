package tests

import "testing"

func Test3IsNot2(t *testing.T) {
	if 3 != 2 {
		t.Error("3 is not 2")
	}
}
