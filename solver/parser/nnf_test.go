package parser

import (
	"fmt"
	"testing"
)

var nnfTestCases = []struct {
	in       string
	expected formula
}{
	{"true", True},
	{"false", False},
}

func TestConstant(t *testing.T) {
	expr, _ := parseExpression("true")
	result := expr.toNNF()
	expected := True
	if !equalFormula(result, expected) {
		t.Errorf("expected %s, was %s", toString(expected), toString(result))
	}
}

func runNnfTests(t *testing.T, testType string) {
	for _, testCase := range nnfTestCases {
		t.Run(testCase.in, func(t *testing.T) {
			expr, _ := parseExpression(testCase.in)
			result := expr.toNNF()
			if !equalFormula(result, testCase.expected) {
				t.Errorf("expected %s, was %s", toString(testCase.expected), toString(result))
			}
		})
	}
}

func equalFormula(f1 formula, f2 formula) bool {
	if fmt.Sprintf("%T", f1) != fmt.Sprintf("%T", f2) {
		return false
	}
	subs1 := f1.subformulas()
	subs2 := f2.subformulas()
	if len(subs1) != len(subs2) {
		return false
	}
	for i, sub1 := range subs1 {
		sub2 := subs2[i]
		if !equalFormula(sub1, sub2) {
			return false
		}
	}
	return true
}
