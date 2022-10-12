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

// retrieveLiterals

func (e Expression) retrieveLiterals() []Literal {
	if e.Implication != nil {
		return e.Implication.retrieveLiterals()
	} else { // e.Unique != nil
		return e.Unique.retrieveLiterals()
	}
}

func (u Unique) retrieveLiterals() []Literal {
	literals := make([]Literal, 0)
	for cur := &u; cur != nil; cur = cur.Next {
		literal := *cur.First
		literals = append(literals, literal)
	}
	return literals
}

func (i Implication) retrieveLiterals() []Literal {
	left := i.Left.retrieveLiterals()
	var right []Literal
	if i.Implication != nil {
		right = i.Implication.retrieveLiterals()

	} else if i.Equivalence != nil {
		right = i.Equivalence.retrieveLiterals()
	}
	return append(left, right...)
}

func (d Disjunction) retrieveLiterals() []Literal {
	operands := make([]Literal, 0)
	for cur := &d; cur != nil; cur = cur.Next {
		operand := cur.Conjunction.retrieveLiterals()
		operands = append(operands, operand...)
	}
	return operands
}

func (c Conjunction) retrieveLiterals() []Literal {
	operands := make([]Literal, 0)
	for cur := &c; cur != nil; cur = cur.Next {
		operand := cur.Unary.retrieveLiterals()
		operands = append(operands, operand...)
	}
	return operands
}

func (u Unary) retrieveLiterals() []Literal {
	if u.Not != "" {
		return u.Unary.retrieveLiterals()
	} else {
		return u.Factor.retrieveLiterals()
	}
}

func (f Factor) retrieveLiterals() []Literal {
	if f.Constant != nil {
		return []Literal{}
	} else if f.Literal != nil {
		return []Literal{*f.Literal}
	} else { // f.SubExpression != nil
		return f.SubExpression.retrieveLiterals()
	}
}
