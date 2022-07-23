package parser

import (
	"reflect"
	"testing"
)

var nnfTestCases = []struct {
	in       string
	expected nnf
}{
	{"true", True},
	{"false", False},
	{"!true", False},
	{"b", lit{name: "b"}},
	{"!a", lit{name: "a", negated: true}},
	{"!(a & b)", or{literalNeg("a"), literalNeg("b")}},
	{"(a & b)", and{literal("a"), literal("b")}},
	{"!(a | !b)", and{literalNeg("a"), literal("b")}},
	{"!a | !b", or{literalNeg("a"), literalNeg("b")}},
	{"a & true", literal("a")},
	{"a & false", False},
	{"a | true", True},
	{"a | false", literal("a")},
	{"((a | b | c) | d | e)", or{literal("a"), literal("b"), literal("c"), literal("d"), literal("e")}},
	{"((a & b & c) & d & e)", and{literal("a"), literal("b"), literal("c"), literal("d"), literal("e")}},
	{"(!a & b) -> !(c & d)", or{literal("a"), literalNeg("b"), literalNeg("c"), literalNeg("d")}},
	{"a = b", and{or{literalNeg("a"), literal("b")}, or{literal("a"), literalNeg("b")}}},
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

func equalFormula(f1 nnf, f2 nnf) bool {
	return reflect.DeepEqual(f1, f2)
}

func literal(name string) lit {
	return lit{name: name, negated: false}
}

func literalNeg(name string) lit {
	return lit{name: name, negated: true}
}
