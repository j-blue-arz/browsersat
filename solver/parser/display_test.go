package parser

import (
	"testing"

	"github.com/alecthomas/participle/v2"
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
}}

func TestParseAndOr(t *testing.T) {
	runTests(t, "and-or")
}

func TestParseFactors(t *testing.T) {
	runTests(t, "factors")
}

func runTests(t *testing.T, testType string) {
	for _, testCase := range testCases[testType] {
		t.Run(testCase.in, func(t *testing.T) {
			p := participle.MustBuild[Expression](participle.UseLookahead(2))
			expr, err := p.ParseString("", testCase.in)
			if err != nil {
				t.Fatal(err)
			}
			result, err := expr.String()
			if err != nil {
				t.Fatal(err)
			}
			if result != testCase.expected {
				t.Fatalf("Parsing result is %s, expected %s", result, testCase.expected)
			}
		})

	}
}
