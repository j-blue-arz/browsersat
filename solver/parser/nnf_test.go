package parser

import (
	"reflect"
	"testing"
)

var nnfTestCases = []struct {
	in       string
	expected formula
}{
	{"true", True},
	{"false", False},
	{"!true", False},
	{"!a", lit{name: "a", negated: true}},
}

func TestAllNnfCases(t *testing.T) {
	for _, testCase := range nnfTestCases {
		t.Run(testCase.in, func(t *testing.T) {
			expr, _ := parseExpression(testCase.in)
			result := expr.toNNF(false)
			if !equalFormula(result, testCase.expected) {
				t.Errorf("expected %s, was %s", toString(testCase.expected), toString(result))
			}
		})
	}
}

func equalFormula(f1 formula, f2 formula) bool {
	return reflect.DeepEqual(f1, f2)
}
