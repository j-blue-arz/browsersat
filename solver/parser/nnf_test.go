package parser

import "testing"

func TestConstant(t *testing.T) {
	expr, _ := parseExpression("true")
	result := expr.nnf()
	if result != True {
		t.Errorf("expected %s, was %s", True, result)
	}
}
