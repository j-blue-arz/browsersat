package maxsat

import (
	"fmt"
	"strings"

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

func AddConstraint(inputConstraint string) error {
	newConstraints := append(inputConstraints, inputConstraint)
	formula, err := parse(newConstraints)
	if err != nil {
		return fmt.Errorf("could not parse constraint %q: %v", inputConstraint, err)
	}
	inputConstraints = newConstraints
	cnf := AsCnf(formula)
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
	formula, _ := parse(inputConstraints)
	cnf := AsCnf(formula)
	cnf.AddUnitLiteral(toFormula(literal, !model[literal]))
	newModel, err := solve(cnf)
	if err != nil {
		return fmt.Errorf("flipping %q leads to UNSAT", literal)
	}
	model = newModel
	return nil
}

func toFormula(literal string, value bool) Formula {
	if value {
		return Var(literal)
	} else {
		return Not(Var(literal))
	}
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

func solve(cnf *Cnf) (map[string]bool, error) {
	newModel := solveMaxsat(cnf)
	if newModel == nil {
		return nil, fmt.Errorf("no model for UNSAT constraints")
	}
	return cnf.TransformModel(newModel), nil
}

func solveMaxsat(cnf *Cnf) []bool {
	for literal, value := range model {
		literalFormula := toFormula(literal, value)
		cnf.AddRelaxableLiteral(literalFormula)
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

func parse(input []string) (Formula, error) {
	constraints := strings.Join(input, "; ")
	reader := strings.NewReader(constraints)
	return Parse(reader)
}
