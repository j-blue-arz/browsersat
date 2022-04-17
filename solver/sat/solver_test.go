package sat

import "testing"

func TestLiteralSat(t *testing.T) {
	result, _ := IsSat("x")
	sat, ok := result["sat"].(bool)
	if !ok || !sat {
		t.Fail()
	}
}

func TestFormulaUnsat(t *testing.T) {
	result, _ := IsSat("x & ^x")
	sat, ok := result["sat"].(bool)
	if !ok || sat {
		t.Fail()
	}
}
