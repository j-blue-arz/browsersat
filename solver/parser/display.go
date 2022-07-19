package parser

import (
	"fmt"
	"strings"
)

func (e Expression) string() (string, error) {
	return e.Implication.string()
}

func (e Implication) string() (string, error) {
	disjunction, err := e.Left.string()
	if err != nil {
		return "", err
	}
	if e.Implication != nil {
		right, err := e.Implication.string()
		if err != nil {
			return "", err
		}
		return "implies(" + disjunction + ", " + right + ")", nil
	} else if e.Equivalence != nil {
		right, err := e.Equivalence.string()
		if err != nil {
			return "", err
		}
		return "eq(" + disjunction + ", " + right + ")", nil
	}
	return disjunction, nil
}

func (e Disjunction) string() (string, error) {
	operands := make([]string, 0)
	for cur := &e; cur != nil; cur = cur.Next {
		operand, err := cur.Conjunction.string()
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

func (e Conjunction) string() (string, error) {
	operands := make([]string, 0)
	for cur := &e; cur != nil; cur = cur.Next {
		operand, err := cur.Unary.string()
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

func (e Unary) string() (string, error) {
	if e.Factor != nil {
		return e.Factor.string()
	}
	result, err := e.Unary.string()
	if err != nil {
		return "", err
	} else {
		return "not(" + result + ")", err
	}
}

func (e Factor) string() (string, error) {
	if e.Constant != nil {
		return e.Constant.string()
	} else if e.Literal != nil {
		return e.Literal.string()
	} else if e.SubExpression != nil {
		return e.SubExpression.string()
	} else {
		return "", fmt.Errorf("factor must be one of constant, literal, or subexpression")
	}
}

func (e Constant) string() (string, error) {
	if *(e.Value) {
		return "TRUE", nil
	} else {
		return "FALSE", nil
	}
}

func (e Literal) string() (string, error) {
	return e.Name, nil
}
