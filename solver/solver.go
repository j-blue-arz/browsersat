package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/crillab/gophersat/bf"
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

func jsonWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		isSat, err := IsSat(args[0].String())
		if err != nil {
			fmt.Printf("IsSat returned error: %s\n", err)
			return err.Error()
		}
		return isSat
	})
	return jsonFunc
}

func main() {
	js.Global().Set("isSat", jsonWrapper())
	<-make(chan bool)
}
