package maxsat

import "github.com/alecthomas/participle/v2"
import "github.com/alecthomas/participle/v2/lexer"

// expression    	→ implication ;
// implication   	→ disj ( "=" disj ) ;
//				  	       | "->" disj )? ;
// disj  			→ conj ( "|" disj )* ;
// conj 			→ unary (  "&"  conj )* ;
// unary      		→ ( "!" ) unary
//             		| factor ;
// factor			→ constant ;
//					| literal
//             		| "(" expression ")" ;
// literal     		→ STRING
// constant			→ true | false

type Expression struct {
	Implication *Implication `parser:"@@"`
	Unique      *Unique      `parser:"|'{' @@ '}'"`
}

type Unique struct {
	First *Literal `parser:"@@"`
	Next  *Unique  `parser:"(','  @@)?"`
}

type Implication struct {
	Left        *Disjunction `parser:"@@"`
	Implication *Disjunction `parser:"( ImplicationOperator  @@"`
	Equivalence *Disjunction `parser:"| EquivalenceOperator  @@ )?"`
}

type Disjunction struct {
	Conjunction *Conjunction `parser:"@@"`
	Next        *Disjunction `parser:"(OrOperator  @@)?"`
}

type Conjunction struct {
	Unary *Unary       `parser:"@@"`
	Next  *Conjunction `parser:"(AndOperator  @@)?"`
}

type Unary struct {
	Not    string  `parser:"( @( NotOperator )"`
	Unary  *Unary  `parser:" @@ )"`
	Factor *Factor `parser:"| @@"`
}

type Factor struct {
	Constant      *Constant    `parser:"@@"`
	Literal       *Literal     `parser:"| @@"`
	SubExpression *Disjunction `parser:"| '(' @@ ')'"`
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type Constant struct {
	Value *Boolean `parser:"@( 'true' | 'false' )"`
}

type Literal struct {
	Name string `parser:"@LiteralName"`
}

var expressionLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "ImplicationOperator", Pattern: `->`},
	{Name: "EquivalenceOperator", Pattern: `=`},
	{Name: "AndOperator", Pattern: `&|\+`},
	{Name: "OrOperator", Pattern: `\||\/`},
	{Name: "NotOperator", Pattern: `!|\-`},
	{Name: "Parentheses", Pattern: `\(|\)`},
	{Name: "UniqueBrace", Pattern: `\{|\}`},
	{Name: "UniqueDelim", Pattern: `,`},
	{Name: "Value", Pattern: `true|false`},
	{Name: "LiteralName", Pattern: `\w+`},
	{Name: "Whitespace", Pattern: `[ \t]+`},
})

func parseExpression(s string) (*Expression, error) {
	expressionParser := participle.MustBuild[Expression](
		participle.UseLookahead(1),
		participle.Lexer(expressionLexer),
		participle.Elide("Whitespace"))
	return expressionParser.ParseString("", s)
}
