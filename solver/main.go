package main

import (
	"fmt"
	"github.com/j-blue-arz/browsersat/solver/sat"
	"syscall/js"
)

func jsonWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		isSat, err := sat.IsSat(args[0].String())
		if err != nil {
			fmt.Printf("IsSat returned error: %s\n", err)
			return err.Error()
		}
		return isSat
	})
	return jsonFunc
}

func main() {
	js.Global().Set("solveFormula", jsonWrapper())
	<-make(chan bool)
}
