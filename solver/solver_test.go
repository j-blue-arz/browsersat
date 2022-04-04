package main

import "testing"

func TestLiteralSat(t *testing.T) {
	result := IsSat("x")
	if !result {
		t.Fail()
	}
}
