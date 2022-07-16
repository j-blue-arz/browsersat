package parser

import "github.com/alecthomas/participle/v2"

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
}

type Implication struct {
	Left        *Disjunction `parser:"@@"`
	Implication *Disjunction `parser:"( '>'  @@"`
	Equivalence *Disjunction `parser:"| '='  @@ )?"`
}

type Disjunction struct {
	Conjunction *Conjunction `parser:"@@"`
	Next        *Disjunction `parser:"('|'  @@)?"`
}

type Conjunction struct {
	Unary *Unary       `parser:"@@"`
	Next  *Conjunction `parser:"('&'  @@)?"`
}

type Unary struct {
	Not    string  `parser:"( @( '!' )"`
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
	Name *string `parser:"@Ident"`
}

func parse(s string) (*Expression, error) {
	p := participle.MustBuild[Expression](participle.UseLookahead(2))
	return p.ParseString("", s)
}
