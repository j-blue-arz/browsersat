package maxsat

import "testing"

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
	val, ok := result[literal].(bool)
	if !ok || !val {
		t.Fail()
	}
}
