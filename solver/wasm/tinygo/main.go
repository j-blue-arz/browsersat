package main

import (
	"fmt"
	"github.com/j-blue-arz/browsersat/solver/maxsat"
)

// export initializeSolver
func initializeSolver() {
	maxsat.Init()
}

// export addConstraint
func addConstraint(constraint string) bool {
	err := maxsat.AddConstraint(constraint)
	if err != nil {
		fmt.Printf("AddConstraint returned error: %s\n", err)
		return false
	}
	return true
}

// export isSat
func isSat() bool {
	return maxsat.IsSat()
}

// export getModel
func getModel() map[string]bool {
	model, err := maxsat.GetModel()
	if err != nil {
		fmt.Printf("GetModel returned error: %s\n", err)
		return nil
	}
	return model
}

// export flipLiteral
func flipLiteral(literal string) bool {
	error := maxsat.FlipLiteral(literal)
	return error == nil
}
