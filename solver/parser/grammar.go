package parser

// expression    	→ equivalence ;
// equivalence   	→ implication ( "=" equivalence )* ;
// implication  	→ disj ( "->" implication )* ;
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
	Equivalence *Equivalence `parser:"@@"`
}

type Equivalence struct {
	Implication *Implication `parser:"@@"`
	Next        *Equivalence `parser:"('='  @@)?"`
}

type Implication struct {
	Disjunction *Disjunction `parser:"@@"`
	Next        *Implication `parser:"('->'  @@)?"`
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
	Constant      *Constant   `parser:"@@"`
	Literal       *Literal    `parser:"| @@"`
	SubExpression *Expression `parser:"| '(' @@ ')'"`
}

type Constant struct {
	Value *bool `parser:"( @'true' | 'false' )"`
}

type Literal struct {
	Name *string `parser:"@Ident"`
}
