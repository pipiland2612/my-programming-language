package evaluator

import (
	"fmt"
	"my-programming-language/ast"
	"my-programming-language/token"
	"strings"
)

// Values

type Value interface {
	String() string
}

type IntVal struct{ Value int }
type BoolVal struct{ Value bool }
type StringVal struct{ Value string }
type UnitVal struct{}
type PairVal struct{ First, Second Value }
type ListVal struct{ Elems []Value }
type InlVal struct{ Value Value }
type InrVal struct{ Value Value }

type ClosureVal struct {
	Param string
	Body  ast.Expr
	Env   *Env
}

func (v IntVal) String() string     { return fmt.Sprintf("%d", v.Value) }
func (v BoolVal) String() string    { return fmt.Sprintf("%t", v.Value) }
func (v StringVal) String() string  { return v.Value }
func (v UnitVal) String() string    { return "()" }
func (v PairVal) String() string    { return fmt.Sprintf("(%s, %s)", v.First.String(), v.Second.String()) }
func (v InlVal) String() string     { return fmt.Sprintf("inl %s", v.Value.String()) }
func (v InrVal) String() string     { return fmt.Sprintf("inr %s", v.Value.String()) }
func (v ClosureVal) String() string { return "<function>" }

func (v ListVal) String() string {
	parts := make([]string, len(v.Elems))
	for i, e := range v.Elems {
		parts[i] = e.String()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

// Environment

type Env struct {
	bindings map[string]Value
	parent   *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{bindings: make(map[string]Value), parent: parent}
}

func (env *Env) Get(name string) (Value, bool) {
	if v, ok := env.bindings[name]; ok {
		return v, true
	}
	if env.parent != nil {
		return env.parent.Get(name)
	}
	return nil, false
}

func (env *Env) Set(name string, v Value) {
	env.bindings[name] = v
}

// Evaluator

type Evaluator struct {
	env    *Env
	output strings.Builder
}

func New() *Evaluator {
	return &Evaluator{env: NewEnv(nil)}
}

func NewWithEnv(env *Env) *Evaluator {
	return &Evaluator{env: env}
}

func (ev *Evaluator) Env() *Env {
	return ev.env
}

func (ev *Evaluator) GetOutput() string {
	return ev.output.String()
}

func (ev *Evaluator) ClearOutput() {
	ev.output.Reset()
}

func errAt(pos token.Pos, msg string) error {
	return fmt.Errorf("%d:%d: runtime error: %s", pos.Line, pos.Column, msg)
}

func (ev *Evaluator) Eval(expr ast.Expr) (Value, error) {
	return ev.eval(expr, ev.env)
}

func (ev *Evaluator) EvalProgram(prog *ast.Program) (Value, error) {
	var lastVal Value = &UnitVal{}
	for _, decl := range prog.Declarations {
		switch d := decl.(type) {
		case *ast.ImportExpr:
			continue // handled externally
		case *ast.LetExpr:
			if d.Body != nil {
				v, err := ev.eval(d, ev.env)
				if err != nil {
					return nil, err
				}
				lastVal = v
			} else {
				v, err := ev.eval(d.Value, ev.env)
				if err != nil {
					return nil, err
				}
				ev.env.Set(d.Name, v)
				lastVal = v
			}
		default:
			v, err := ev.eval(decl, ev.env)
			if err != nil {
				return nil, err
			}
			lastVal = v
		}
	}
	return lastVal, nil
}

func (ev *Evaluator) eval(expr ast.Expr, env *Env) (Value, error) {
	switch e := expr.(type) {
	case *ast.IntLit:
		return &IntVal{Value: e.Value}, nil

	case *ast.BoolLit:
		return &BoolVal{Value: e.Value}, nil

	case *ast.StringLit:
		return &StringVal{Value: e.Value}, nil

	case *ast.UnitLit:
		return &UnitVal{}, nil

	case *ast.Var:
		v, ok := env.Get(e.Name)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("undefined variable '%s'", e.Name))
		}
		return v, nil

	case *ast.BinOp:
		return ev.evalBinOp(e, env)

	case *ast.UnaryOp:
		return ev.evalUnaryOp(e, env)

	case *ast.IfExpr:
		condV, err := ev.eval(e.Cond, env)
		if err != nil {
			return nil, err
		}
		b := condV.(*BoolVal)
		if b.Value {
			return ev.eval(e.Then, env)
		}
		return ev.eval(e.Else, env)

	case *ast.LetExpr:
		valV, err := ev.eval(e.Value, env)
		if err != nil {
			return nil, err
		}
		if e.Body == nil {
			env.Set(e.Name, valV)
			return valV, nil
		}
		newEnv := NewEnv(env)
		newEnv.Set(e.Name, valV)
		return ev.eval(e.Body, newEnv)

	case *ast.FnExpr:
		return &ClosureVal{Param: e.Param, Body: e.Body, Env: env}, nil

	case *ast.AppExpr:
		funcV, err := ev.eval(e.Func, env)
		if err != nil {
			return nil, err
		}
		argV, err := ev.eval(e.Arg, env)
		if err != nil {
			return nil, err
		}
		closure, ok := funcV.(*ClosureVal)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("cannot apply non-function value: %s", funcV.String()))
		}
		appEnv := NewEnv(closure.Env)
		appEnv.Set(closure.Param, argV)
		return ev.eval(closure.Body, appEnv)

	case *ast.PairExpr:
		fstV, err := ev.eval(e.First, env)
		if err != nil {
			return nil, err
		}
		sndV, err := ev.eval(e.Second, env)
		if err != nil {
			return nil, err
		}
		return &PairVal{First: fstV, Second: sndV}, nil

	case *ast.FstExpr:
		v, err := ev.eval(e.Expr, env)
		if err != nil {
			return nil, err
		}
		p := v.(*PairVal)
		return p.First, nil

	case *ast.SndExpr:
		v, err := ev.eval(e.Expr, env)
		if err != nil {
			return nil, err
		}
		p := v.(*PairVal)
		return p.Second, nil

	case *ast.ListExpr:
		elems := make([]Value, len(e.Elems))
		for i, elem := range e.Elems {
			v, err := ev.eval(elem, env)
			if err != nil {
				return nil, err
			}
			elems[i] = v
		}
		return &ListVal{Elems: elems}, nil

	case *ast.ConsExpr:
		headV, err := ev.eval(e.Head, env)
		if err != nil {
			return nil, err
		}
		tailV, err := ev.eval(e.Tail, env)
		if err != nil {
			return nil, err
		}
		list := tailV.(*ListVal)
		newElems := make([]Value, 0, len(list.Elems)+1)
		newElems = append(newElems, headV)
		newElems = append(newElems, list.Elems...)
		return &ListVal{Elems: newElems}, nil

	case *ast.InlExpr:
		v, err := ev.eval(e.Expr, env)
		if err != nil {
			return nil, err
		}
		return &InlVal{Value: v}, nil

	case *ast.InrExpr:
		v, err := ev.eval(e.Expr, env)
		if err != nil {
			return nil, err
		}
		return &InrVal{Value: v}, nil

	case *ast.CaseExpr:
		return ev.evalCase(e, env)

	case *ast.FixExpr:
		return ev.evalFix(e, env)

	case *ast.PrintExpr:
		v, err := ev.eval(e.Expr, env)
		if err != nil {
			return nil, err
		}
		if e.Newline {
			ev.output.WriteString(v.String() + "\n")
			fmt.Println(v.String())
		} else {
			ev.output.WriteString(v.String())
			fmt.Print(v.String())
		}
		return &UnitVal{}, nil

	case *ast.ImportExpr:
		return &UnitVal{}, nil

	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (ev *Evaluator) evalBinOp(e *ast.BinOp, env *Env) (Value, error) {
	leftV, err := ev.eval(e.Left, env)
	if err != nil {
		return nil, err
	}
	rightV, err := ev.eval(e.Right, env)
	if err != nil {
		return nil, err
	}

	switch e.Op {
	case "+":
		if ls, ok := leftV.(*StringVal); ok {
			rs := rightV.(*StringVal)
			return &StringVal{Value: ls.Value + rs.Value}, nil
		}
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &IntVal{Value: l.Value + r.Value}, nil
	case "-":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &IntVal{Value: l.Value - r.Value}, nil
	case "*":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &IntVal{Value: l.Value * r.Value}, nil
	case "/":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		if r.Value == 0 {
			return nil, errAt(e.Pos, "division by zero")
		}
		return &IntVal{Value: l.Value / r.Value}, nil
	case "%":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		if r.Value == 0 {
			return nil, errAt(e.Pos, "modulo by zero")
		}
		return &IntVal{Value: l.Value % r.Value}, nil
	case "==":
		return &BoolVal{Value: valuesEqual(leftV, rightV)}, nil
	case "!=":
		return &BoolVal{Value: !valuesEqual(leftV, rightV)}, nil
	case "<":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &BoolVal{Value: l.Value < r.Value}, nil
	case ">":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &BoolVal{Value: l.Value > r.Value}, nil
	case "<=":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &BoolVal{Value: l.Value <= r.Value}, nil
	case ">=":
		l, r := leftV.(*IntVal), rightV.(*IntVal)
		return &BoolVal{Value: l.Value >= r.Value}, nil
	case "&&":
		l, r := leftV.(*BoolVal), rightV.(*BoolVal)
		return &BoolVal{Value: l.Value && r.Value}, nil
	case "||":
		l, r := leftV.(*BoolVal), rightV.(*BoolVal)
		return &BoolVal{Value: l.Value || r.Value}, nil
	}
	return nil, errAt(e.Pos, fmt.Sprintf("unknown operator '%s'", e.Op))
}

func (ev *Evaluator) evalUnaryOp(e *ast.UnaryOp, env *Env) (Value, error) {
	v, err := ev.eval(e.Expr, env)
	if err != nil {
		return nil, err
	}
	switch e.Op {
	case "!":
		b := v.(*BoolVal)
		return &BoolVal{Value: !b.Value}, nil
	case "-":
		i := v.(*IntVal)
		return &IntVal{Value: -i.Value}, nil
	}
	return nil, errAt(e.Pos, fmt.Sprintf("unknown unary operator '%s'", e.Op))
}

func (ev *Evaluator) evalCase(e *ast.CaseExpr, env *Env) (Value, error) {
	scrutV, err := ev.eval(e.Scrutinee, env)
	if err != nil {
		return nil, err
	}

	for _, branch := range e.Branches {
		branchEnv := NewEnv(env)

		switch p := branch.Pattern.(type) {
		case *ast.NilPattern:
			if list, ok := scrutV.(*ListVal); ok && len(list.Elems) == 0 {
				return ev.eval(branch.Body, branchEnv)
			}
		case *ast.ConsPattern:
			if list, ok := scrutV.(*ListVal); ok && len(list.Elems) > 0 {
				branchEnv.Set(p.Head, list.Elems[0])
				branchEnv.Set(p.Tail, &ListVal{Elems: list.Elems[1:]})
				return ev.eval(branch.Body, branchEnv)
			}
		case *ast.InlPattern:
			if inl, ok := scrutV.(*InlVal); ok {
				branchEnv.Set(p.Name, inl.Value)
				return ev.eval(branch.Body, branchEnv)
			}
		case *ast.InrPattern:
			if inr, ok := scrutV.(*InrVal); ok {
				branchEnv.Set(p.Name, inr.Value)
				return ev.eval(branch.Body, branchEnv)
			}
		}
	}

	return nil, errAt(e.Pos, "no matching case branch")
}

func (ev *Evaluator) evalFix(e *ast.FixExpr, env *Env) (Value, error) {
	v, err := ev.eval(e.Expr, env)
	if err != nil {
		return nil, err
	}
	closure := v.(*ClosureVal)

	// fix(g) = g(fix(g))
	// For fix (fn (f: A) => body): evaluate body with f = result, where result = fix(g)
	// Use mutation to tie the knot
	fixEnv := NewEnv(closure.Env)
	result, err := ev.eval(closure.Body, fixEnv)
	if err != nil {
		return nil, err
	}
	// Now set the self-reference so f points to the result
	fixEnv.Set(closure.Param, result)
	return result, nil
}

func valuesEqual(a, b Value) bool {
	switch a := a.(type) {
	case *IntVal:
		if b, ok := b.(*IntVal); ok {
			return a.Value == b.Value
		}
	case *BoolVal:
		if b, ok := b.(*BoolVal); ok {
			return a.Value == b.Value
		}
	case *StringVal:
		if b, ok := b.(*StringVal); ok {
			return a.Value == b.Value
		}
	case *UnitVal:
		_, ok := b.(*UnitVal)
		return ok
	case *PairVal:
		if b, ok := b.(*PairVal); ok {
			return valuesEqual(a.First, b.First) && valuesEqual(a.Second, b.Second)
		}
	case *ListVal:
		if b, ok := b.(*ListVal); ok {
			if len(a.Elems) != len(b.Elems) {
				return false
			}
			for i := range a.Elems {
				if !valuesEqual(a.Elems[i], b.Elems[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}
