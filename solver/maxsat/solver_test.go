package maxsat

import (
	"fmt"
	"testing"
)

func TestLiteralSat(t *testing.T) {
	Init()
	err := AddConstraint("x")
	if err != nil {
		t.Fail()
	}
	assertSat(t)
	assertTrue("x", t)
}

func TestFormulaUnsat(t *testing.T) {
	Init()
	err := AddConstraint("x & ^x")
	if err != nil {
		t.Fail()
	}
	assertUnsat(t)
}

func TestAddConstraintMultipleTimes(t *testing.T) {
	Init()
	AddConstraint("a | ^b")
	assertSat(t)
	AddConstraint("b | ^c")
	assertSat(t)
	AddConstraint("c | a")
	assertSat(t)
	AddConstraint("^a")
	assertUnsat(t)
}

func TestInit(t *testing.T) {
	Init()
	AddConstraint("a")
	assertSat(t)
	AddConstraint("^a")
	assertUnsat(t)
	Init()
	AddConstraint("a")
	assertSat(t)
}

func TestFlipLiteral(t *testing.T) {
	Init()
	AddConstraint("a | b | c")
	model, _ := GetModel()
	FlipLiteral("b")
	newModel, _ := GetModel()
	if model["b"] == newModel["b"] {
		t.Errorf("FlipLiteral should flip model[\"b\"], but remained %t", model["b"])
	}
}

func assertSat(t *testing.T) {
	sat := IsSat()
	if !sat {
		t.Fail()
	}
}

func assertUnsat(t *testing.T) {
	sat := IsSat()
	if sat {
		t.Fail()
	}
}

func assertTrue(literal string, t *testing.T) {
	result, _ := GetModel()
	val := result[literal]
	if !val {
		t.Fail()
	}
}

func assertEq(a, b int, t *testing.T) {
	if a != b {
		t.Fail()
	}
}
