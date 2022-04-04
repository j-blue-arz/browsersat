package main

import (
	"fmt"
	"github.com/crillab/gophersat/bf"
	"strings"
)

func IsSat(inputFormula string) (bool, error) {
	reader := strings.NewReader(inputFormula)
	formula, err := bf.Parse(reader)
	if err != nil {
		return false, fmt.Errorf("could not parse formula %q: %v", inputFormula, err)
	}
	model := bf.Solve(formula)
	return model != nil, nil
}
