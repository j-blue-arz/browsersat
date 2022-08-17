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

func toLit(literal string, value bool) lit {
	return lit{name: literal, negated: !value}
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
	if s.Solve() != solver.Sat {
		return nil
	}
	m := s.Model()
	return m
}
