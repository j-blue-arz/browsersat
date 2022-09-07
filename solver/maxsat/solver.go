package maxsat

import (
	"fmt"

	"github.com/j-blue-arz/tiny-gophersat/solver"
)

var inputConstraints []string
var model map[string]bool
var isSat bool

func Init() {
	inputConstraints = make([]string, 0)
	model = make(map[string]bool)
	isSat = false
}

func ValidateConstraint(constraint string) error {
	_, err := parseExpression(constraint)
	if err != nil {
		return err
	}
	return nil
}

func AddConstraint(inputConstraint string) error {
	newConstraints := append(inputConstraints, inputConstraint)
	cnf, err := parseToCnf(newConstraints)
	if err != nil {
		return fmt.Errorf("could not parse constraint %q: %v", inputConstraint, err)
	}
	inputConstraints = newConstraints
	newModel, err := solve(cnf)

	if err != nil {
		isSat = false
		model = make(map[string]bool)
	} else {
		isSat = true
		model = newModel
	}
	return nil
}

func FlipLiteral(literal string) error {
	if !isSat {
		return fmt.Errorf("constraints are not satisfiable")
	}
	if _, ok := model[literal]; !ok {
		return fmt.Errorf("literal %q not contained in model", literal)
	}
	cnf, _ := parseToCnf(inputConstraints)
	cnf.addUnitLiteral(toLit(literal, !model[literal]))
	newModel, err := solve(cnf)
	if err != nil {
		return fmt.Errorf("flipping %q leads to UNSAT", literal)
	}
	model = newModel
	return nil
}

func IsSat() bool {
	return isSat
}

func GetModel() (map[string]bool, error) {
	if isSat {
		return model, nil
	} else {
		return nil, fmt.Errorf("no model for UNSAT constraints")
	}
}

func Evaluate(constraint string) (bool, error) {
	expr, err := parseExpression(constraint)
	if err != nil {
		return false, fmt.Errorf("not a valid expression")
	}

	literals := expr.retrieveLiterals()
	for _, lit := range literals {
		if _, ok := model[lit.Name]; !ok {
			return false, fmt.Errorf("expression contains unknown literals")
		}
	}

	return expr.evaluate(model)
}

func toLit(literal string, value bool) lit {
	return lit{name: literal, negated: !value}
}

func solve(cnf *cnf) (map[string]bool, error) {
	newModel := solveMaxsat(cnf)
	if newModel == nil {
		return nil, fmt.Errorf("no model for UNSAT constraints")
	}
	return cnf.transformModel(newModel), nil
}

func solveMaxsat(cnf *cnf) []bool {
	for literal, value := range model {
		cnf.addRelaxableLiteral(toLit(literal, value))
	}
	weights := make([]int, len(cnf.relaxLits))
	relaxLits := make([]solver.Lit, len(cnf.relaxLits))

	for i, lit := range cnf.relaxLits {
		weights[i] = 1
		relaxLits[i] = solver.IntToLit(int32(lit))
	}

	problem := solver.ParseSlice(cnf.clauses)
	if len(cnf.relaxLits) > 0 {
		problem.SetCostFunc(relaxLits, weights)
	}
	s := solver.New(problem)
	res := s.Optimal(nil, nil)
	if res.Status != solver.Sat {
		return nil
	}
	return res.Model
}

// retrieveLiterals

func (e Expression) retrieveLiterals() []Literal {
	return e.Implication.retrieveLiterals()
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

// evaluate

func (e Expression) evaluate(model map[string]bool) (bool, error) {
	return e.Implication.evaluate(model)
}

func (i Implication) evaluate(model map[string]bool) (bool, error) {
	disjunction, err := i.Left.evaluate(model)
	if err != nil {
		return false, err
	}
	if i.Implication != nil {
		right, err := i.Implication.evaluate(model)
		if err != nil {
			return false, err
		}
		return !disjunction || right, nil
	} else if i.Equivalence != nil {
		right, err := i.Equivalence.evaluate(model)
		if err != nil {
			return false, err
		}
		return disjunction == right, nil
	}
	return disjunction, nil
}

func (d Disjunction) evaluate(model map[string]bool) (bool, error) {
	for cur := &d; cur != nil; cur = cur.Next {
		operand, err := cur.Conjunction.evaluate(model)
		if err != nil {
			return false, err
		}
		if operand {
			return true, nil
		}
	}
	return false, nil
}

func (c Conjunction) evaluate(model map[string]bool) (bool, error) {
	for cur := &c; cur != nil; cur = cur.Next {
		operand, err := cur.Unary.evaluate(model)
		if err != nil {
			return false, err
		}
		if !operand {
			return false, nil
		}
	}
	return true, nil
}

func (u Unary) evaluate(model map[string]bool) (bool, error) {
	if u.Not != "" {
		result, err := u.Unary.evaluate(model)
		if err != nil {
			return false, err
		} else {
			return !result, nil
		}
	} else {
		return u.Factor.evaluate(model)
	}
}

func (f Factor) evaluate(model map[string]bool) (bool, error) {
	if f.Constant != nil {
		return f.Constant.evaluate(model)
	} else if f.Literal != nil {
		return f.Literal.evaluate(model)
	} else if f.SubExpression != nil {
		return f.SubExpression.evaluate(model)
	} else {
		return false, fmt.Errorf("factor must be one of constant, literal, or subexpression")
	}
}

func (c Constant) evaluate(model map[string]bool) (bool, error) {
	if *(c.Value) {
		return true, nil
	} else {
		return false, nil
	}
}

func (l Literal) evaluate(model map[string]bool) (bool, error) {
	if val, ok := model[l.Name]; !ok {
		return false, fmt.Errorf("expression contains unknown literals")
	} else {
		return val, nil
	}
}
