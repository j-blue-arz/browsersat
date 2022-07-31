// This source file transforms a set of constraints into a cnf to be used by the sat solver
// Large parts were copied from the bf package of the gophersat project

package maxsat

import "fmt"

type cnf struct {
	vars      vars
	clauses   [][]int
	relaxLits []int
}

func parseToCnf(constraints []string) (*cnf, error) {
	vars := vars{all: make(map[string]int), problem: make(map[string]int)}
	var clauses [][]int
	for _, constraint := range constraints {
		expression, err := parseExpression(constraint)
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, cnfRec(expression.toNNF(false), &vars)...)
	}
	return &cnf{vars: vars, clauses: clauses}, nil
}

func (cnf *cnf) transformModel(model []bool) map[string]bool {
	vars := make(map[string]bool)
	for name, idx := range cnf.vars.problem {
		vars[name] = model[idx-1]
	}
	return vars
}

func (cnf *cnf) addUnitLiteral(l lit) error {
	if !cnf.vars.known(l) {
		return fmt.Errorf("cannot add unit literal %s, because it is not a known variable", l.name)
	}
	clause := []int{cnf.vars.litValue(l)}
	cnf.clauses = append(cnf.clauses, clause)
	return nil
}

func (cnf *cnf) addRelaxableLiteral(l lit) error {
	if !cnf.vars.known(l) {
		return fmt.Errorf("cannot add relaxable literal %s, because it is not a known variable", l.name)
	}
	relaxLit := cnf.vars.dummy()
	clause := []int{cnf.vars.litValue(l), relaxLit}
	cnf.clauses = append(cnf.clauses, clause)
	cnf.relaxLits = append(cnf.relaxLits, relaxLit)
	return nil
}

func cnfRec(f nnf, vars *vars) [][]int {
	switch f := f.(type) {
	case lit:
		return [][]int{{vars.litValue(f)}}
	case and:
		var res [][]int
		for _, sub := range f {
			res = append(res, cnfRec(sub, vars)...)
		}
		return res
	case or:
		var res [][]int
		var lits []int
		for _, sub := range f {
			switch sub := sub.(type) {
			case lit:
				lits = append(lits, vars.litValue(sub))
			case and:
				d := vars.dummy()
				lits = append(lits, d)
				for _, sub2 := range sub {
					cnf := cnfRec(sub2, vars)
					for i := range cnf {
						cnf[i] = append(cnf[i], -d)
					}
					res = append(res, cnf...)
				}
			default:
				panic("unexpected or in or")
			}
		}
		res = append(res, lits)
		return res
	case trueConst:
		return [][]int{}
	case falseConst:
		return [][]int{{}}
	default:
		panic("invalid NNF formula")
	}
}

type vars struct {
	all     map[string]int
	problem map[string]int
}

func (vars *vars) known(l lit) bool {
	_, ok := vars.all[l.name]
	return ok
}

func (vars *vars) litValue(l lit) int {
	val, ok := vars.all[l.name]
	if !ok {
		val = len(vars.all) + 1
		vars.all[l.name] = val
		vars.problem[l.name] = val
	}
	if l.negated {
		return -val
	}
	return val
}

func (vars *vars) dummy() int {
	val := len(vars.all) + 1
	vars.all[fmt.Sprintf("dummy-%d", val)] = val
	return val
}
