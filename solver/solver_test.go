package main

import "testing"

func TestLiteralSat(t *testing.T) {
	result, _ := IsSat("x")
	if !result {
		t.Fail()
	}
}

func TestFormulaUnsat(t *testing.T) {
	result, _ := IsSat("x & ^x")
	if result {
		t.Fail()
	}
}
