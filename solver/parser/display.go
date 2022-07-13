package parser

import (
	"fmt"
	"strings"
)

type Formula interface {
	String() (string, error)
}

func (e Expression) String() (string, error) {
	return e.Equivalence.String()
}

func (e Equivalence) String() (string, error) {
	return e.Implication.String()
}

func (e Implication) String() (string, error) {
	return e.Disjunction.String()
}

func (e Disjunction) String() (string, error) {
	operands := make([]string, 0)
	for cur := &e; cur != nil; cur = cur.Next {
		operand, err := cur.Conjunction.String()
		if err != nil {
			return "", err
		}
		operands = append(operands, operand)
	}
	if len(operands) > 1 {
		return "or(" + strings.Join(operands, ", ") + ")", nil
	} else if len(operands) == 1 {
		return operands[0], nil
	} else {
		return "", fmt.Errorf("expected at least one operand, got 0")
	}
}

func (e Conjunction) String() (string, error) {
	operands := make([]string, 0)
	for cur := &e; cur != nil; cur = cur.Next {
		operand, err := cur.Unary.String()
		if err != nil {
			return "", err
		}
		operands = append(operands, operand)
	}
	if len(operands) > 1 {
		return "and(" + strings.Join(operands, ", ") + ")", nil
	} else if len(operands) == 1 {
		return operands[0], nil
	} else {
		return "", fmt.Errorf("expected at least one operand, got 0")
	}
}

func (e Unary) String() (string, error) {
	if e.Factor != nil {
		return e.Factor.String()
	}
	result, err := e.Unary.String()
	if err != nil {
		return "!" + result, err
	} else {
		return "", err
	}

}

func (e Factor) String() (string, error) {
	return e.Literal.String()
}

func (e Constant) String() (string, error) {
	if *e.Value {
		return "true", nil
	} else {
		return "false", nil
	}
}

func (e Literal) String() (string, error) {
	return *e.Name, nil
}
