package maxsat

import (
	"fmt"
	"strings"

	"github.com/crillab/gophersat/solver"
)

var inputConstraints []string
var model map[string]bool
var isSat bool

func Init() {
	inputConstraints = make([]string, 0)
	model = make(map[string]bool)
	isSat = false
}

func AddConstraint(inputConstraint string) error {
	newConstraints := append(inputConstraints, inputConstraint)
	formula, err := parse(newConstraints)
	if err != nil {
		return fmt.Errorf("could not parse constraint %q: %v", inputConstraint, err)
	}
	inputConstraints = newConstraints
	cnf := AsCnf(formula)
	updateModel(cnf)
	return nil
}

func FlipLiteral(literal string) error {
	if !isSat {
		return fmt.Errorf("constraints are not satisfiable")
	}
	if _, ok := model[literal]; !ok {
		return fmt.Errorf("literal %q not contained in model", literal)
	}
	formula, _ := parse(inputConstraints)
	cnf := AsCnf(formula)
	if model[literal] {
		cnf.AddUnitLiteral(Not(Var(literal)))
	} else {
		cnf.AddUnitLiteral(Var(literal))
	}
	newModel := solve(cnf)
	if newModel == nil {
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

func updateModel(cnf *Cnf) {
	model = solve(cnf)

	if model != nil {
		isSat = true
	} else {
		isSat = false
	}
}

func solve(cnf *Cnf) map[string]bool {
	model := solveMaxsat(cnf)
	if model == nil {
		return nil
	}
	return cnf.TransformModel(model)
}

func solveMaxsat(cnf *Cnf) []bool {
	pb := solver.ParseSlice(cnf.clauses)
	s := solver.New(pb)
	if s.Solve() != solver.Sat {
		return nil
	}
	m := s.Model()
	return m
}

func parse(input []string) (Formula, error) {
	constraints := strings.Join(input, "; ")
	reader := strings.NewReader(constraints)
	return Parse(reader)
}
