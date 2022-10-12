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

func TestImplication(t *testing.T) {
	Init()
	err := AddConstraint("a -> b")
	if err != nil {
		t.Fail()
	}
	assertSat(t)
}

func TestFormulaUnsat(t *testing.T) {
	Init()
	err := AddConstraint("x & !x")
	if err != nil {
		t.Fail()
	}
	assertUnsat(t)
}

func TestAddConstraintMultipleTimes(t *testing.T) {
	Init()
	AddConstraint("a | !b")
	assertSat(t)
	AddConstraint("b | !c")
	assertSat(t)
	AddConstraint("c | a")
	assertSat(t)
	AddConstraint("!a")
	assertUnsat(t)
}

func TestAddConstraintWithUnique(t *testing.T) {
	Init()
	AddConstraint("{a, b, c}")
	assertSat(t)
	AddConstraint("{c, d}")
	assertSat(t)
	AddConstraint("d")
	assertSat(t)
	AddConstraint("!a & !b")
	assertUnsat(t)
}

func TestAddConstraintParserError(t *testing.T) {
	cases := []string{"!", "!!", "a||b", "&7&b"}
	for _, input := range cases {
		Init()
		AddConstraint("a | !b")
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
	AddConstraint("!a")
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
	AddConstraint("a | (b & c & d) | e")
	model, _ := GetModel()
	if model["b"] {
		FlipLiteral("b")
		FlipLiteral("c")
		FlipLiteral("d")
	}
	assertFalse("b", t)
	assertFalse("c", t)
	assertFalse("d", t)
	for i := 1; i <= 20; i++ {
		model, _ := GetModel()
		if model["a"] {
			FlipLiteral("a")
		} else if model["e"] {
			FlipLiteral("e")
		} else {
			t.Fatalf("b, c, and d are expected to remain false, but they became true in iteration %d.", i)
		}
	}
}

func TestValidateConstraint(t *testing.T) {
	_, err := ValidateConstraint("a | (b | c")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestValidateConstraintEmptyUnique(t *testing.T) {
	_, err := ValidateConstraint("{}")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestValidateConstraintNegatedLiteralInUnique(t *testing.T) {
	_, err := ValidateConstraint("{!a, !b}")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestValidateConstraintReturnsCanonicalForm(t *testing.T) {
	str, err := ValidateConstraint("a + -b")
	if err != nil {
		t.Fatalf("unexpected error")
	}
	if str != "and(a, not(b))" {
		t.Errorf("expected %s, got %s", "and(a, not(b))", str)
	}
}

func TestEvaluateInvalidFormula(t *testing.T) {
	Init()
	AddConstraint("a & b & !c & !d")
	_, err := Evaluate("a | (b")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestEvaluateFormulaUnknownLiterals(t *testing.T) {
	Init()
	AddConstraint("a & b & !c & !d")
	_, err := Evaluate("a | f")
	if err == nil {
		t.Errorf("expected error, but got none")
	}
}

func TestEvaluateOr(t *testing.T) {
	Init()
	AddConstraint("a & b & !c & !d")
	expectEvaluateTrue(t, "a | c")
}

func TestEvaluateAndSubformula(t *testing.T) {
	Init()
	AddConstraint("a & b & !c & !d")
	expectEvaluateFalse(t, "a & (!a | c)")
}

func TestEvaluateConstantAndImplication(t *testing.T) {
	Init()
	AddConstraint("a & b & !c & !d")
	expectEvaluateTrue(t, "false -> (!a | c)")
}

func TestEvaluateUniqueConstraint(t *testing.T) {
	Init()
	AddConstraint("a & b & !c & !d")
	expectEvaluateTrue(t, "{a}")
	expectEvaluateTrue(t, "{b, c, d}")
	expectEvaluateFalse(t, "{a, b}")
	expectEvaluateFalse(t, "{c, d}")
}

func expectEvaluateTrue(t *testing.T, constraint string) {
	value, err := Evaluate(constraint)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !value {
		t.Errorf("expected true, was false")
	}
}

func expectEvaluateFalse(t *testing.T, constraint string) {
	value, err := Evaluate(constraint)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if value {
		t.Errorf("expected false, was true")
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
