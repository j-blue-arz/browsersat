package maxsat

import (
	"encoding/json"
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

func TestDigitAsLiteral(t *testing.T) {
	Init()
	AddConstraint("7 | 8")
	AddConstraint("^8")
	assertSat(t)
	assertTrue("7", t)
	assertFalse("8", t)
}

func TestAddConstraintParserError(t *testing.T) {
	cases := []string{"^", "^^", "a||b", "&7&b"}
	for _, input := range cases {
		Init()
		AddConstraint("a | ^b")
		err := AddConstraint(input)
		if err == nil {
			t.Errorf("Expected an parsig error for input '%s', was nil", input)
		}
		assertSat(t)
		AddConstraint("b")
		assertSat(t)
		assertTrue("a", t)
	}
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

func TestFlipLiteralMinimizesModelChanges(t *testing.T) {
	Init()
	AddConstraint("a | b | c")
	model, _ := GetModel()
	literalToFlip := "a"
	if !allTrue(model) {
		literalToFlip, _ = getAnyFalse(model)
	}
	FlipLiteral(literalToFlip)
	newModel, _ := GetModel()

	d := diff(model, newModel)
	if d > 1 {
		t.Errorf("when flipping %s, models %s and %s should only differ in one literal, but was %d", literalToFlip, toStr(model), toStr(newModel), d)
	}
}

func diff(model1, model2 map[string]bool) int {
	sum := 0
	for lit := range model1 {
		if model1[lit] != model2[lit] {
			sum += 1
		}
	}
	return sum
}

func getAnyFalse(model map[string]bool) (string, error) {
	for lit, v := range model {
		if !v {
			return lit, nil
		}
	}
	return "", fmt.Errorf("all values are true")
}

func allTrue(model map[string]bool) bool {
	for _, v := range model {
		if !v {
			return false
		}
	}
	return true
}

func assertSat(t *testing.T) {
	sat := IsSat()
	if !sat {
		t.Errorf("expected SAT, but was UNSAT")
	}
}

func assertUnsat(t *testing.T) {
	sat := IsSat()
	if sat {
		t.Errorf("expected UNSAT, but was SAT")
	}
}

func assertTrue(literal string, t *testing.T) {
	assertLiteralValue(literal, true, t)
}

func assertFalse(literal string, t *testing.T) {
	assertLiteralValue(literal, false, t)
}

func assertLiteralValue(literal string, expected bool, t *testing.T) {
	result, _ := GetModel()
	actual := result[literal]
	if actual != expected {
		t.Errorf("expected %q to be %s, but was %t", literal, boolToStr(expected), actual)
	}
}

func boolToStr(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func toStr(model map[string]bool) string {
	mJson, err := json.Marshal(model)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(mJson)
}
