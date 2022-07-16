package parser

import (
	"testing"
)

var testCases = map[string][]struct {
	in       string
	expected string
}{"and-or": {
	{"a|b", "or(a, b)"},
	{"a&b", "and(a, b)"},
	{"a&b&c", "and(a, b, c)"},
	{"a|b|c", "or(a, b, c)"},
	{"a|b&c|d", "or(a, and(b, c), d)"},
	{"a&b|c&d", "or(and(a, b), and(c, d))"},
	{"a&b|c|d&e", "or(and(a, b), c, and(d, e))"},
}, "factors": {
	{"foo", "foo"},
	{"!foo", "not(foo)"},
	{"!foo | true", "or(not(foo), true)"},
	{"false & true", "and(false, true)"},
}, "subexpressions": {
	{"(foo)", "foo"},
	{"(a | b)", "or(a, b)"},
	{"(a & (b & (c & d)) & e)", "and(a, and(b, and(c, d)), e)"},
	{"((a|b)&(c|(d)))", "and(or(a, b), or(c, d))"},
	{"a&((b|(c&d|e))&f|g)", "and(a, or(and(or(b, or(and(c, d), e)), f), g))"},
}, "implies-equiv": {
	{"a>b", "implies(a, b)"},
	{"a=b", "eq(a, b)"},
	{"a&b=c|d", "eq(and(a, b), or(c, d))"},
	{"a&b>c|d", "implies(and(a, b), or(c, d))"},
}, "errors": {
	{"a=b=c", ""},
	{"a>b>c", ""},
	{"a(|b)|c", ""},
	{"a|(b|c", ""},
	{"a|()|c", ""},
}}

func TestParseAndOr(t *testing.T) {
	runTests(t, "and-or")
}

func TestParseFactors(t *testing.T) {
	runTests(t, "factors")
}

func TestParseSubexpressions(t *testing.T) {
	runTests(t, "subexpressions")
}

func TestParseImplications(t *testing.T) {
	runTests(t, "implies-equiv")
}

func TestParseErrors(t *testing.T) {
	for _, testCase := range testCases["errors"] {
		t.Run(testCase.in, func(t *testing.T) {
			_, err := parse(testCase.in)
			if err == nil {
				t.Fatalf("Parser should error, but did not")
			}
		})
	}
}

func runTests(t *testing.T, testType string) {
	for _, testCase := range testCases[testType] {
		t.Run(testCase.in, func(t *testing.T) {
			expr, err := parse(testCase.in)
			if err != nil {
				t.Fatal(err)
			}
			result, err := expr.string()
			if err != nil {
				t.Fatal(err)
			}
			if result != testCase.expected {
				t.Fatalf("Parsing result is %s, expected %s", result, testCase.expected)
			}
		})
	}
}
