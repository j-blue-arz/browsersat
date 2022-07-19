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
func (l lit) symbol() string         { return l.name }

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

func (e Expression) toNNF() formula {
	return e.Implication.toNNF()
}

func (i Implication) toNNF() formula {
	return i.Left.toNNF()
}

func (d Disjunction) toNNF() formula {
	return d.Conjunction.toNNF()
}

func (c Conjunction) toNNF() formula {
	return c.Unary.toNNF()
}

func (u Unary) toNNF() formula {
	return u.Factor.toNNF()
}

func (f Factor) toNNF() formula {
	return f.Constant.toNNF()
}

func (c Constant) toNNF() formula {
	if *(c.Value) {
		return True
	} else {
		return False
	}
}
