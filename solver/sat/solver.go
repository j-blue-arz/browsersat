package sat

import (
	"fmt"
	"strings"

	"github.com/crillab/gophersat/bf"
)

func IsSat(inputFormula string) (map[string]interface{}, error) {
	reader := strings.NewReader(inputFormula)
	formula, err := bf.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("could not parse formula %q: %v", inputFormula, err)
	}
	model := bf.Solve(formula)
	result := make(map[string]interface{})

	if model != nil {
		result["sat"] = true
		result["model"] = convertModel(model)
	} else {
		result["sat"] = false
	}

	return result, nil
}

func convertModel(model map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	for lit, value := range model {
		result[lit] = value
	}
	return result
}
