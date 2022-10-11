package maxsat

import (
	"fmt"
	"strings"
)

func (e Expression) string() (string, error) {
	if e.Implication != nil {
		return e.Implication.string()
	} else if e.Unique != nil {
		return e.Unique.string()
	} else {
		return "", fmt.Errorf("expected either a boolean formula or a max one constraint")
	}
}

func (u Unique) string() (string, error) {
	literals := make([]string, 0)
	for cur := &u; cur != nil; cur = cur.Next {
		literal, err := cur.First.string()
		if err != nil {
			return "", err
		}
		literals = append(literals, literal)
	}
	if len(literals) > 1 {
		return "unique(" + strings.Join(literals, ", ") + ")", nil
	} else if len(literals) == 1 {
		return literals[0], nil
	} else {
		return "", fmt.Errorf("expected at least one operand, got 0")
	}
}

func (i Implication) string() (string, error) {
	disjunction, err := i.Left.string()
	if err != nil {
		return "", err
	}
	if i.Implication != nil {
		right, err := i.Implication.string()
		if err != nil {
			return "", err
		}
		return "implies(" + disjunction + ", " + right + ")", nil
	} else if i.Equivalence != nil {
		right, err := i.Equivalence.string()
		if err != nil {
			return "", err
		}
		return "eq(" + disjunction + ", " + right + ")", nil
	}
	return disjunction, nil
}

func (d Disjunction) string() (string, error) {
	operands := make([]string, 0)
	for cur := &d; cur != nil; cur = cur.Next {
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

func (c Conjunction) string() (string, error) {
	operands := make([]string, 0)
	for cur := &c; cur != nil; cur = cur.Next {
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

func (u Unary) string() (string, error) {
	if u.Not != "" {
		result, err := u.Unary.string()
		if err != nil {
			return "", err
		} else {
			return "not(" + result + ")", err
		}
	} else {
		return u.Factor.string()
	}
}

func (f Factor) string() (string, error) {
	if f.Constant != nil {
		return f.Constant.string()
	} else if f.Literal != nil {
		return f.Literal.string()
	} else if f.SubExpression != nil {
		return f.SubExpression.string()
	} else {
		return "", fmt.Errorf("factor must be one of constant, literal, or subexpression")
	}
}

func (c Constant) string() (string, error) {
	if *(c.Value) {
		return "TRUE", nil
	} else {
		return "FALSE", nil
	}
}

func (l Literal) string() (string, error) {
	return l.Name, nil
}
