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
	updateModel(formula)
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

func updateModel(formula Formula) {
	model = solve(formula)

	if model != nil {
		isSat = true
	} else {
		isSat = false
	}
}

func solve(formula Formula) map[string]bool {
	cnf := AsCnf(formula)
	pb := solver.ParseSlice(cnf.clauses)
	s := solver.New(pb)
	if s.Solve() != solver.Sat {
		return nil
	}
	m := s.Model()
	vars := make(map[string]bool)
	for v, idx := range cnf.vars.pb {
		vars[v.name] = m[idx-1]
	}
	return vars
}

func parse(input []string) (Formula, error) {
	constraints := strings.Join(input, "; ")
	reader := strings.NewReader(constraints)
	return Parse(reader)
}
