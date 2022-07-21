package parser

import "strings"

type formula interface {
	nnf() formula
	subformulas() []formula
	symbol() string
}

type trueConst struct{}
type falseConst struct{}

var True trueConst = trueConst{}
var False falseConst = falseConst{}

func (t trueConst) nnf() formula           { return t }
func (t trueConst) subformulas() []formula { return make([]formula, 0) }
func (t trueConst) symbol() string         { return "true" }

func (f falseConst) nnf() formula           { return f }
func (f falseConst) subformulas() []formula { return make([]formula, 0) }
func (f falseConst) symbol() string         { return "false" }

type and []formula

func (a and) nnf() formula           { return a }
func (a and) subformulas() []formula { return a }
func (a and) symbol() string         { return "and" }

type or []formula

func (o or) nnf() formula           { return o }
func (o or) subformulas() []formula { return o }
func (o or) symbol() string         { return "or" }

type not []formula

func (n not) nnf() formula           { return n }
func (n not) subformulas() []formula { return n }
func (n not) symbol() string         { return "not" }

type lit struct {
	name    string
	negated bool
}

func (l lit) nnf() formula           { return l }
func (l lit) subformulas() []formula { return make([]formula, 0) }
func (l lit) symbol() string {
	if l.negated {
		return "!" + l.name
	} else {
		return l.name
	}
}

func toString(f formula) string {
	subformulas := f.subformulas()
	if len(subformulas) == 0 {
		return f.symbol()
	} else {
		operands := make([]string, 0, len(subformulas))
		for _, subformula := range subformulas {
			operands = append(operands, toString(subformula))
		}
		return f.symbol() + "(" + strings.Join(operands, ", ") + ")"
	}
}

func (e Expression) toNNF(negated bool) formula {
	return e.Implication.toNNF(negated)
}

func (i Implication) toNNF(negated bool) formula {
	return i.Left.toNNF(negated)
}

func (d Disjunction) toNNF(negated bool) formula {
	return d.Conjunction.toNNF(negated)
}

func (c Conjunction) toNNF(negated bool) formula {
	operands := make([]formula, 0)
	for cur := &c; cur != nil; cur = cur.Next {
		operand := cur.Unary.toNNF(negated)
		operands = append(operands, operand)
	}
	if len(operands) == 1 {
		return operands[0]
	} else if len(operands) > 1 {
		if negated {
			return or(operands)
		} else {
			return and(operands)
		}
	} else { // empty conjunction
		return True
	}
}

func (u Unary) toNNF(negated bool) formula {
	if u.Not != "" {
		return u.Unary.toNNF(!negated)
	} else {
		return u.Factor.toNNF(negated)
	}
}

func (f Factor) toNNF(negated bool) formula {
	if f.Constant != nil {
		return f.Constant.toNNF(negated)
	} else if f.Literal != nil {
		return f.Literal.toNNF(negated)
	} else {
		return f.SubExpression.toNNF(negated)
	}
}

func (c Constant) toNNF(negated bool) formula {
	if bool(*(c.Value)) == negated {
		return False
	} else {
		return True
	}
}

func (l Literal) toNNF(negated bool) formula {
	return lit{name: l.Name, negated: negated}
}
