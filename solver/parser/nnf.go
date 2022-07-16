package parser

type formula interface {
	nnf() formula
}

type trueConst struct{}
type falseConst struct{}

/* type and []formula
type or []formula
type not formula
*/
var True trueConst = trueConst{}
var False falseConst = falseConst{}

func (t trueConst) nnf() formula  { return t }
func (f falseConst) nnf() formula { return f }

func (e Expression) nnf() formula {
	return e.Implication.nnf()
}

func (i Implication) nnf() formula {
	return i.Left.nnf()
}

func (d Disjunction) nnf() formula {
	return d.Conjunction.nnf()
}

func (c Conjunction) nnf() formula {
	return c.Unary.nnf()
}

func (u Unary) nnf() formula {
	return u.Factor.nnf()
}

func (f Factor) nnf() formula {
	return f.Constant.nnf()
}

func (e Constant) nnf() formula {
	if *(e.Value) {
		return True
	} else {
		return False
	}
}
