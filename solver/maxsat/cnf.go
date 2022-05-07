// This source file transforms a set of constraints into a cnf to be used by the sat solver
// Large parts were copied from the bf package of the gophersat project

package maxsat

import (
	"fmt"
	"io"
	"text/scanner"
)

// A Formula is any kind of boolean formula, not necessarily in CNF.
type Formula interface {
	nnf() Formula
}

type trueConst struct{}

var True Formula = trueConst{}

func (t trueConst) nnf() Formula { return t }

type falseConst struct{}

var False Formula = falseConst{}

func (f falseConst) nnf() Formula { return f }

type variable struct {
	name  string
	dummy bool
}

func pbVar(name string) variable {
	return variable{name: name, dummy: false}
}

func dummyVar(name string) variable {
	return variable{name: name, dummy: true}
}

func Var(name string) Formula {
	return pbVar(name)
}

func (v variable) nnf() Formula {
	return lit{signed: false, v: v}
}

type lit struct {
	v      variable
	signed bool
}

func (l lit) nnf() Formula {
	return l
}

type not [1]Formula

func Not(f Formula) Formula {
	return not{f}
}

func (n not) nnf() Formula {
	switch f := n[0].(type) {
	case variable:
		l := f.nnf().(lit)
		l.signed = true
		return l
	case lit:
		f.signed = !f.signed
		return f
	case not:
		return f[0].nnf()
	case and:
		subs := make([]Formula, len(f))
		for i, sub := range f {
			subs[i] = not{sub}.nnf()
		}
		return or(subs).nnf()
	case or:
		subs := make([]Formula, len(f))
		for i, sub := range f {
			subs[i] = not{sub}.nnf()
		}
		return and(subs).nnf()
	case trueConst:
		return False
	case falseConst:
		return True
	default:
		panic("invalid formula type")
	}
}

type and []Formula

func And(subs ...Formula) Formula {
	return and(subs)
}

func (a and) nnf() Formula {
	var res and
	for _, s := range a {
		nnf := s.nnf()
		switch nnf := nnf.(type) {
		case and: // Simplify: "and"s in the "and" get to the higher level
			res = append(res, nnf...)
		case trueConst: // True is ignored
		case falseConst:
			return False
		default:
			res = append(res, nnf)
		}
	}
	if len(res) == 1 {
		return res[0]
	}
	if len(res) == 0 {
		return False
	}
	return res
}

type or []Formula

func Or(subs ...Formula) Formula {
	return or(subs)
}

func (o or) nnf() Formula {
	var res or
	for _, s := range o {
		nnf := s.nnf()
		switch nnf := nnf.(type) {
		case or: // Simplify: "or"s in the "or" get to the higher level
			res = append(res, nnf...)
		case falseConst: // False is ignored
		case trueConst:
			return True
		default:
			res = append(res, nnf)
		}
	}
	if len(res) == 1 {
		return res[0]
	}
	if len(res) == 0 {
		return True
	}
	return res
}

func Implies(f1, f2 Formula) Formula {
	return or{not{f1}, f2}
}

func Eq(f1, f2 Formula) Formula {
	return and{or{not{f1}, f2}, or{f1, not{f2}}}
}

func Xor(f1, f2 Formula) Formula {
	return and{or{not{f1}, not{f2}}, or{f1, f2}}
}

type vars struct {
	all map[variable]int // all vars, including those created when converting the formula
	pb  map[variable]int // Only the vars that appeared orinigally in the problem
}

func (vars *vars) litValue(l lit) int {
	val, ok := vars.all[l.v]
	if !ok {
		val = len(vars.all) + 1
		vars.all[l.v] = val
		vars.pb[l.v] = val
	}
	if l.signed {
		return -val
	}
	return val
}

func (vars *vars) dummy() int {
	val := len(vars.all) + 1
	vars.all[dummyVar(fmt.Sprintf("dummy-%d", val))] = val
	return val
}

type Cnf struct {
	vars      vars
	clauses   [][]int
	relaxLits []int
}

// asCnf returns a CNF representation of the given formula.
func AsCnf(f Formula) *Cnf {
	vars := vars{all: make(map[variable]int), pb: make(map[variable]int)}
	clauses := cnfRec(f.nnf(), &vars)
	return &Cnf{vars: vars, clauses: clauses}
}

func cnfRec(f Formula, vars *vars) [][]int {
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
					cnf[0] = append(cnf[0], -d)
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

type parser struct {
	s     scanner.Scanner
	eof   bool
	token string
}

func Parse(r io.Reader) (Formula, error) {
	var s scanner.Scanner
	s.Init(r)
	p := parser{s: s}
	p.scan()
	f, err := p.parseClause()
	if err != nil {
		return f, err
	}
	if !p.eof {
		return nil, fmt.Errorf("expected EOF, found %q at %v", p.token, p.s.Pos())
	}
	return f, nil
}

func isOperator(token string) bool {
	return token == "=" || token == "->" || token == "|" || token == "&" || token == ";"
}

func (p *parser) scan() {
	p.eof = p.eof || (p.s.Scan() == scanner.EOF)
	p.token = p.s.TokenText()
}

func (p *parser) parseClause() (f Formula, err error) {
	if isOperator(p.token) {
		return nil, fmt.Errorf("unexpected token %q at %s", p.token, p.s.Pos())
	}
	f, err = p.parseEquiv()
	if err != nil {
		return nil, err
	}
	if p.eof {
		return f, nil
	}
	if p.token == ";" {
		p.scan()
		if p.eof {
			return f, nil
		}
		f2, err := p.parseClause()
		if err != nil {
			return nil, err
		}
		return And(f, f2), nil
	}
	return f, nil
}

func (cnf *Cnf) TransformModel(model []bool) map[string]bool {
	vars := make(map[string]bool)
	for v, idx := range cnf.vars.pb {
		vars[v.name] = model[idx-1]
	}
	return vars
}

// Only supports Var("x") and Not(Var("x")), and only if "x" is known in vars
// TODO: make this enforced at compile time (introduce type Literal)
func (cnf *Cnf) AddUnitLiteral(f Formula) {
	clause := cnfRec(f.nnf(), &cnf.vars)[0]
	cnf.clauses = append(cnf.clauses, clause)
}

func (cnf *Cnf) AddRelaxableLiteral(f Formula) {
	clause := cnfRec(f.nnf(), &cnf.vars)[0]
	relaxLit := cnf.vars.dummy()
	clause = append(clause, relaxLit)
	cnf.clauses = append(cnf.clauses, clause)
	cnf.relaxLits = append(cnf.relaxLits, relaxLit)
}

func (p *parser) parseEquiv() (f Formula, err error) {
	if p.eof {
		return nil, fmt.Errorf("at position %v, expected expression, found EOF", p.s.Pos())
	}
	if isOperator(p.token) {
		return nil, fmt.Errorf("unexpected token %q at %s", p.token, p.s.Pos())
	}
	f, err = p.parseImplies()
	if err != nil {
		return nil, err
	}
	if p.eof {
		return f, nil
	}
	if p.token == "=" {
		p.scan()
		if p.eof {
			return nil, fmt.Errorf("unexpected EOF")
		}
		f2, err := p.parseEquiv()
		if err != nil {
			return nil, err
		}
		return Eq(f, f2), nil
	}
	return f, nil
}

func (p *parser) parseImplies() (f Formula, err error) {
	f, err = p.parseOr()
	if err != nil {
		return nil, err
	}
	if p.eof {
		return f, nil
	}
	if p.token == "-" {
		p.scan()
		if p.eof {
			return nil, fmt.Errorf("unexpected EOF")
		}
		if p.token != ">" {
			return nil, fmt.Errorf("invalid token %q at %v", "-"+p.token, p.s.Pos())
		}
		p.scan()
		if p.eof {
			return nil, fmt.Errorf("unexpected EOF")
		}
		f2, err := p.parseImplies()
		if err != nil {
			return nil, err
		}
		return Implies(f, f2), nil
	}
	return f, nil
}

func (p *parser) parseOr() (f Formula, err error) {
	f, err = p.parseAnd()
	if err != nil {
		return nil, err
	}
	if p.eof {
		return f, nil
	}
	if p.token == "|" {
		p.scan()
		if p.eof {
			return nil, fmt.Errorf("unexpected EOF")
		}
		f2, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		return Or(f, f2), nil
	}
	return f, nil
}

func (p *parser) parseAnd() (f Formula, err error) {
	f, err = p.parseNot()
	if err != nil {
		return nil, err
	}
	if p.eof {
		return f, nil
	}
	if p.token == "&" {
		p.scan()
		if p.eof {
			return nil, fmt.Errorf("unexpected EOF")
		}
		f2, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		return And(f, f2), nil
	}
	return f, nil
}

func (p *parser) parseNot() (f Formula, err error) {
	if isOperator(p.token) {
		return nil, fmt.Errorf("unexpected token %q at %s", p.token, p.s.Pos())
	}
	if p.token == "^" {
		p.scan()
		if p.eof {
			return nil, fmt.Errorf("unexpected EOF")
		}
		f, err = p.parseNot()
		if err != nil {
			return nil, err
		}
		return Not(f), nil
	}
	f, err = p.parseBasic()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (p *parser) parseBasic() (f Formula, err error) {
	if isOperator(p.token) || p.token == ")" {
		return nil, fmt.Errorf("unexpected token %q at %s", p.token, p.s.Pos())
	}
	if p.token == "(" {
		p.scan()
		f, err = p.parseEquiv()
		if err != nil {
			return nil, err
		}
		if p.eof {
			return nil, fmt.Errorf("expected closing parenthesis, found EOF at %s", p.s.Pos())
		}
		if p.token != ")" {
			return nil, fmt.Errorf("expected closing parenthesis, found %q at %s", p.token, p.s.Pos())
		}
		p.scan()
		return f, nil
	}
	defer p.scan()
	return Var(p.token), nil
}
