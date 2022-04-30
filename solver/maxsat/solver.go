package maxsat

import (
	"fmt"
	"strings"

	"github.com/crillab/gophersat/bf"
)

var inputConstraints []string
var jsonModel map[string]interface{}
var isSat bool

func Init() {
	inputConstraints = make([]string, 0)
	jsonModel = make(map[string]interface{})
	isSat = false
}

func AddConstraint(inputConstraint string) error {
	inputConstraints = append(inputConstraints, inputConstraint)
	err := updateModel()
	return err
}

func IsSat() bool {
	return isSat
}

func GetModel() (map[string]interface{}, error) {
	if isSat {
		return jsonModel, nil
	} else {
		return nil, fmt.Errorf("no model for UNSAT constraints")
	}

}

func updateModel() error {
	constraints := strings.Join(inputConstraints, "; ")
	reader := strings.NewReader(constraints)
	formula, err := bf.Parse(reader)
	if err != nil {
		return fmt.Errorf("could not parse constraints %q: %v", constraints, err)
	}
	model := bf.Solve(formula)
	jsonModel = make(map[string]interface{})

	if model != nil {
		jsonModel = convertModel(model)
		isSat = true
	} else {
		isSat = false
	}
	return nil
}

func convertModel(model map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	for lit, value := range model {
		result[lit] = value
	}
	return result
}
