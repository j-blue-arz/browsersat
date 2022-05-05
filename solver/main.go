package main

import (
	"fmt"
	"github.com/j-blue-arz/browsersat/solver/maxsat"
	"syscall/js"
)

func initializeSolver(this js.Value, args []js.Value) interface{} {
	maxsat.Init()
	return true
}

func addConstraint(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "Invalid no of arguments passed"
	}
	err := maxsat.AddConstraint(args[0].String())
	if err != nil {
		fmt.Printf("AddConstraint returned error: %s\n", err)
		return err.Error()
	}
	return true
}

func isSat(this js.Value, args []js.Value) interface{} {
	return maxsat.IsSat()
}

func getModel(this js.Value, args []js.Value) interface{} {
	model, err := maxsat.GetModel()
	if err != nil {
		fmt.Printf("GetModel returned error: %s\n", err)
		return err.Error()
	}
	return convertModel(model)
}

func flipLiteral(this js.Value, args []js.Value) interface{} {
	error := maxsat.FlipLiteral(args[0].String())
	if error != nil {
		return false
	}
	return true
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
	export["addConstraint"] = js.FuncOf(addConstraint)
	export["isSat"] = js.FuncOf(isSat)
	export["getModel"] = js.FuncOf(getModel)
	export["flipLiteral"] = js.FuncOf(flipLiteral)
	js.Global().Set("satsolver", export)

	<-make(chan bool)
}
