package main

import "testing"

func TestLiteralSat(t *testing.T) {
	result, _ := IsSat("x")
	if !result {
		t.Fail()
	}
}
