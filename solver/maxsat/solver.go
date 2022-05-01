package maxsat

import (
	"fmt"
	"strings"

	"github.com/crillab/gophersat/bf"
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

func updateModel(formula bf.Formula) {
	model = bf.Solve(formula)

	if model != nil {
		isSat = true
	} else {
		isSat = false
	}
}

func parse(input []string) (bf.Formula, error) {
	constraints := strings.Join(input, "; ")
	reader := strings.NewReader(constraints)
	return bf.Parse(reader)
}
