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
	// negated == false, because Implication can only be top-level operator in grammar
	if i.Implication != nil {
		left := i.Left.toNNF(true)
		right := i.Implication.toNNF(false)
		return makeOr([]nnf{left, right})
	} else {
		return i.Left.toNNF(false)
	}
}

func (d Disjunction) toNNF(negated bool) nnf {
	operands := make([]nnf, 0)
	for cur := &d; cur != nil; cur = cur.Next {
		operand := cur.Conjunction.toNNF(negated)
		operands = append(operands, operand)
	}
	if negated {
		return makeAnd(operands)
	} else {
		return makeOr(operands)
	}
}

func (c Conjunction) toNNF(negated bool) nnf {
	operands := make([]nnf, 0)
	for cur := &c; cur != nil; cur = cur.Next {
		operand := cur.Unary.toNNF(negated)
		operands = append(operands, operand)
	}
	if negated {
		return makeOr(operands)
	} else {
		return makeAnd(operands)
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

func makeOr(operands []nnf) nnf {
	var result or
	for _, f := range operands {
		switch f := f.(type) {
		case or:
			result = append(result, f...)
		case falseConst:
		case trueConst:
			return True
		default:
			result = append(result, f)
		}
	}
	if len(result) == 0 {
		return False
	} else if len(result) == 1 {
		return result[0]
	} else {
		return result
	}
}

func makeAnd(operands []nnf) nnf {
	var result and
	for _, f := range operands {
		switch f := f.(type) {
		case and:
			result = append(result, f...)
		case falseConst:
			return False
		case trueConst:
		default:
			result = append(result, f)
		}
	}
	if len(result) == 0 {
		return True
	} else if len(result) == 1 {
		return result[0]
	} else {
		return result
	}
}
