package parser

import "strings"

type nnf interface {
	subformulas() []nnf
	symbol() string
}

type trueConst struct{}
type falseConst struct{}

var True trueConst = trueConst{}
var False falseConst = falseConst{}

func (t trueConst) subformulas() []nnf { return make([]nnf, 0) }
func (t trueConst) symbol() string     { return "true" }

func (f falseConst) subformulas() []nnf { return make([]nnf, 0) }
func (f falseConst) symbol() string     { return "false" }

type and []nnf

func (a and) subformulas() []nnf { return a }
func (a and) symbol() string     { return "and" }

type or []nnf

func (o or) subformulas() []nnf { return o }
func (o or) symbol() string     { return "or" }

type lit struct {
	name    string
	negated bool
}

func (l lit) subformulas() []nnf { return make([]nnf, 0) }
func (l lit) symbol() string {
	if l.negated {
		return "!" + l.name
	} else {
		return l.name
	}
}

func toString(f nnf) string {
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

func (e Expression) toNNF(negated bool) nnf {
	return e.Implication.toNNF(negated)
}

func (i Implication) toNNF(negated bool) nnf {
	return i.Left.toNNF(negated)
}

func (d Disjunction) toNNF(negated bool) nnf {
	return d.Conjunction.toNNF(negated)
}

func (c Conjunction) toNNF(negated bool) nnf {
	operands := make([]nnf, 0)
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

func (u Unary) toNNF(negated bool) nnf {
	if u.Not != "" {
		return u.Unary.toNNF(!negated)
	} else {
		return u.Factor.toNNF(negated)
	}
}

func (f Factor) toNNF(negated bool) nnf {
	if f.Constant != nil {
		return f.Constant.toNNF(negated)
	} else if f.Literal != nil {
		return f.Literal.toNNF(negated)
	} else {
		return f.SubExpression.toNNF(negated)
	}
}

func (c Constant) toNNF(negated bool) nnf {
	if bool(*(c.Value)) == negated {
		return False
	} else {
		return True
	}
}

func (l Literal) toNNF(negated bool) nnf {
	return lit{name: l.Name, negated: negated}
}
