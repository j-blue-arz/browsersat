//go:build wasm

package main

import (
	"github.com/j-blue-arz/browsersat/solver/maxsat"
	"syscall/js"
)

func initializeSolver(this js.Value, args []js.Value) interface{} {
	maxsat.Init()
	return true
}

func addConstraint(this js.Value, args []js.Value) interface{} {
	err := maxsat.AddConstraint(args[0].String())
	if err != nil {
		return false
	}
	return true
}

func validateConstraint(this js.Value, args []js.Value) interface{} {
	_, err := maxsat.ValidateConstraint(args[0].String())
	if err != nil {
		return err.Error()
	}
	return "VALID"
}

func isSat(this js.Value, args []js.Value) interface{} {
	return maxsat.IsSat()
}

func getModel(this js.Value, args []js.Value) interface{} {
	model, err := maxsat.GetModel()
	if err != nil {
		return err.Error()
	}
	return convertModel(model)
}

func flipLiteral(this js.Value, args []js.Value) interface{} {
	err := maxsat.FlipLiteral(args[0].String())
	if err != nil {
		return false
	}
	return true
}

func evaluate(this js.Value, args []js.Value) interface{} {
	value, err := maxsat.Evaluate(args[0].String())
	if err != nil {
		return "INVALID"
	}
	return value
}

func convertModel(model map[string]bool) map[string]interface{} {
	result := make(map[string]interface{})
	for lit, value := range model {
		result[lit] = value
	}
	return result
}

func main() {
	export := make(map[string]interface{})
	export["initializeSolver"] = js.FuncOf(initializeSolver)
	export["validateConstraint"] = js.FuncOf(validateConstraint)
	export["addConstraint"] = js.FuncOf(addConstraint)
	export["isSat"] = js.FuncOf(isSat)
	export["getModel"] = js.FuncOf(getModel)
	export["flipLiteral"] = js.FuncOf(flipLiteral)
	export["evaluate"] = js.FuncOf(evaluate)
	js.Global().Set("satsolver", export)

	<-make(chan bool)
}
