package main

import (
	"fmt"
	"github.com/j-blue-arz/browsersat/solver/sat"
	"syscall/js"
)

func solveFormula(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		return "Invalid no of arguments passed"
	}
	result, err := sat.IsSat(args[0].String())
	if err != nil {
		fmt.Printf("IsSat returned error: %s\n", err)
		return err.Error()
	}

	return result
}

func main() {
	js.Global().Set("solveFormula", js.FuncOf(solveFormula))
	<-make(chan bool)
}
