package maxsat

import (
	"testing"
)

var displayTestCases = map[string][]struct {
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
	{"!foo | true", "or(not(foo), TRUE)"},
	{"false & true", "and(FALSE, TRUE)"},
}, "subexpressions": {
	{"(foo)", "foo"},
	{"(a | b)", "or(a, b)"},
	{"(!(a & b) | !(!c))", "or(not(and(a, b)), not(not(c)))"},
	{"(a & (b & (c & d)) & e)", "and(a, and(b, and(c, d)), e)"},
	{"((a|b)&(c|(d)))", "and(or(a, b), or(c, d))"},
	{"a&((b|(c&d|e))&f|g)", "and(a, or(and(or(b, or(and(c, d), e)), f), g))"},
}, "implies-equiv": {
	{"a->b", "implies(a, b)"},
	{"a=b", "eq(a, b)"},
	{"a&b=c|d", "eq(and(a, b), or(c, d))"},
	{"a&b->c|d", "implies(and(a, b), or(c, d))"},
}, "errors": {
	{"a=b=c", ""},
	{"a->b->c", ""},
	{"a(|b)|c", ""},
	{"a|(b|c", ""},
	{"a|()|c", ""},
}}

func TestParseAndOr(t *testing.T) {
	runDisplayTests(t, "and-or")
}

func TestParseFactors(t *testing.T) {
	runDisplayTests(t, "factors")
}

func TestParseSubexpressions(t *testing.T) {
	runDisplayTests(t, "subexpressions")
}

func TestParseImplications(t *testing.T) {
	runDisplayTests(t, "implies-equiv")
}

func TestParseErrors(t *testing.T) {
	for _, testCase := range displayTestCases["errors"] {
		t.Run(testCase.in, func(t *testing.T) {
			_, err := parseExpression(testCase.in)
			if err == nil {
				t.Fatalf("Parser should error, but did not")
			}
		})
	}
}

func runDisplayTests(t *testing.T, testType string) {
	for _, testCase := range displayTestCases[testType] {
		t.Run(testCase.in, func(t *testing.T) {
			expr, err := parseExpression(testCase.in)
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
