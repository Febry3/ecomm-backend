package tests

import "testing"

func Test2Is2(t *testing.T) {
	if 2 != 2 {
		t.Error("2 is not 2")
	}
}
