package checker

import (
	"fmt"
	"my-programming-language/ast"
	"my-programming-language/token"
)

type TypeEnv struct {
	bindings map[string]ast.Type
	parent   *TypeEnv
}

func NewTypeEnv(parent *TypeEnv) *TypeEnv {
	return &TypeEnv{bindings: make(map[string]ast.Type), parent: parent}
}

func (env *TypeEnv) Get(name string) (ast.Type, bool) {
	if t, ok := env.bindings[name]; ok {
		return t, true
	}
	if env.parent != nil {
		return env.parent.Get(name)
	}
	return nil, false
}

func (env *TypeEnv) Set(name string, t ast.Type) {
	env.bindings[name] = t
}

// ADTDef stores the definition of an algebraic data type
type ADTDef struct {
	Name     string
	Variants []ast.VariantDef
}

type Checker struct {
	env  *TypeEnv
	ADTs map[string]*ADTDef // type name -> definition
}

func New() *Checker {
	return &Checker{env: NewTypeEnv(nil), ADTs: make(map[string]*ADTDef)}
}

func NewWithEnv(env *TypeEnv) *Checker {
	return &Checker{env: env, ADTs: make(map[string]*ADTDef)}
}

func (c *Checker) Env() *TypeEnv {
	return c.env
}

func errAt(pos token.Pos, msg string) error {
	return fmt.Errorf("%d:%d: type error: %s", pos.Line, pos.Column, msg)
}

func typesEqual(a, b ast.Type) bool {
	switch a := a.(type) {
	case *ast.IntType:
		_, ok := b.(*ast.IntType)
		return ok
	case *ast.BoolType:
		_, ok := b.(*ast.BoolType)
		return ok
	case *ast.StringType:
		_, ok := b.(*ast.StringType)
		return ok
	case *ast.UnitType:
		_, ok := b.(*ast.UnitType)
		return ok
	case *ast.FuncType:
		b, ok := b.(*ast.FuncType)
		return ok && typesEqual(a.Param, b.Param) && typesEqual(a.Return, b.Return)
	case *ast.PairType:
		b, ok := b.(*ast.PairType)
		return ok && typesEqual(a.First, b.First) && typesEqual(a.Second, b.Second)
	case *ast.ListType:
		b, ok := b.(*ast.ListType)
		return ok && typesEqual(a.Elem, b.Elem)
	case *ast.SumType:
		b, ok := b.(*ast.SumType)
		return ok && typesEqual(a.Left, b.Left) && typesEqual(a.Right, b.Right)
	case *ast.TupleType:
		b, ok := b.(*ast.TupleType)
		if !ok || len(a.Elems) != len(b.Elems) {
			return false
		}
		for i := range a.Elems {
			if !typesEqual(a.Elems[i], b.Elems[i]) {
				return false
			}
		}
		return true
	case *ast.RecordType:
		b, ok := b.(*ast.RecordType)
		if !ok || len(a.Fields) != len(b.Fields) {
			return false
		}
		for i := range a.Fields {
			if a.Fields[i].Name != b.Fields[i].Name || !typesEqual(a.Fields[i].Type, b.Fields[i].Type) {
				return false
			}
		}
		return true
	case *ast.ADTType:
		b, ok := b.(*ast.ADTType)
		return ok && a.Name == b.Name
	}
	return false
}

func (c *Checker) Check(expr ast.Expr) (ast.Type, error) {
	return c.check(expr, c.env)
}

func (c *Checker) CheckProgram(prog *ast.Program) error {
	for _, decl := range prog.Declarations {
		switch d := decl.(type) {
		case *ast.ImportExpr:
			// Imports are handled externally
			continue
		case *ast.TypeDecl:
			if err := c.checkTypeDecl(d); err != nil {
				return err
			}
			continue
		case *ast.LetExpr:
			if d.Body != nil {
				// This is a let-in expression at top level, just type check it
				if _, err := c.check(d, c.env); err != nil {
					return err
				}
				continue
			}
			t, err := c.check(d.Value, c.env)
			if err != nil {
				return err
			}
			if d.TypeAnn != nil {
				if !typesEqual(t, d.TypeAnn) {
					return errAt(d.Pos, fmt.Sprintf("declaration '%s' has type %s but annotated as %s",
						d.Name, t.String(), d.TypeAnn.String()))
				}
			}
			c.env.Set(d.Name, t)
		default:
			// Bare expression at top level
			if _, err := c.check(decl, c.env); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Checker) check(expr ast.Expr, env *TypeEnv) (ast.Type, error) {
	switch e := expr.(type) {
	case *ast.IntLit:
		return &ast.IntType{}, nil

	case *ast.BoolLit:
		return &ast.BoolType{}, nil

	case *ast.StringLit:
		return &ast.StringType{}, nil

	case *ast.UnitLit:
		return &ast.UnitType{}, nil

	case *ast.Var:
		t, ok := env.Get(e.Name)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("undefined variable '%s'", e.Name))
		}
		return t, nil

	case *ast.BinOp:
		return c.checkBinOp(e, env)

	case *ast.UnaryOp:
		return c.checkUnaryOp(e, env)

	case *ast.IfExpr:
		condT, err := c.check(e.Cond, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(condT, &ast.BoolType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("if condition must be Bool, got %s", condT.String()))
		}
		thenT, err := c.check(e.Then, env)
		if err != nil {
			return nil, err
		}
		elseT, err := c.check(e.Else, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(thenT, elseT) {
			return nil, errAt(e.Pos, fmt.Sprintf("if branches must have same type: then has %s, else has %s",
				thenT.String(), elseT.String()))
		}
		return thenT, nil

	case *ast.LetExpr:
		valT, err := c.check(e.Value, env)
		if err != nil {
			return nil, err
		}
		if e.TypeAnn != nil {
			if !typesEqual(valT, e.TypeAnn) {
				return nil, errAt(e.Pos, fmt.Sprintf("let '%s' has type %s but annotated as %s",
					e.Name, valT.String(), e.TypeAnn.String()))
			}
		}
		if e.Body == nil {
			// Top-level declaration
			env.Set(e.Name, valT)
			return valT, nil
		}
		newEnv := NewTypeEnv(env)
		newEnv.Set(e.Name, valT)
		return c.check(e.Body, newEnv)

	case *ast.FnExpr:
		newEnv := NewTypeEnv(env)
		newEnv.Set(e.Param, e.ParamType)
		bodyT, err := c.check(e.Body, newEnv)
		if err != nil {
			return nil, err
		}
		return &ast.FuncType{Param: e.ParamType, Return: bodyT}, nil

	case *ast.AppExpr:
		funcT, err := c.check(e.Func, env)
		if err != nil {
			return nil, err
		}
		ft, ok := funcT.(*ast.FuncType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("cannot apply non-function type %s", funcT.String()))
		}
		argT, err := c.check(e.Arg, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(ft.Param, argT) {
			return nil, errAt(e.Pos, fmt.Sprintf("function expects %s but got %s",
				ft.Param.String(), argT.String()))
		}
		return ft.Return, nil

	case *ast.PairExpr:
		fstT, err := c.check(e.First, env)
		if err != nil {
			return nil, err
		}
		sndT, err := c.check(e.Second, env)
		if err != nil {
			return nil, err
		}
		return &ast.PairType{First: fstT, Second: sndT}, nil

	case *ast.FstExpr:
		t, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		pt, ok := t.(*ast.PairType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("fst expects a pair, got %s", t.String()))
		}
		return pt.First, nil

	case *ast.SndExpr:
		t, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		pt, ok := t.(*ast.PairType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("snd expects a pair, got %s", t.String()))
		}
		return pt.Second, nil

	case *ast.ListExpr:
		if len(e.Elems) == 0 {
			if e.ElemType != nil {
				return &ast.ListType{Elem: e.ElemType}, nil
			}
			return nil, errAt(e.Pos, "empty list needs a type annotation: [] : [Type]")
		}
		elemT, err := c.check(e.Elems[0], env)
		if err != nil {
			return nil, err
		}
		for i := 1; i < len(e.Elems); i++ {
			t, err := c.check(e.Elems[i], env)
			if err != nil {
				return nil, err
			}
			if !typesEqual(elemT, t) {
				return nil, errAt(e.Elems[i].GetPos(), fmt.Sprintf("list elements must have same type: expected %s, got %s",
					elemT.String(), t.String()))
			}
		}
		return &ast.ListType{Elem: elemT}, nil

	case *ast.ConsExpr:
		headT, err := c.check(e.Head, env)
		if err != nil {
			return nil, err
		}
		tailT, err := c.check(e.Tail, env)
		if err != nil {
			return nil, err
		}
		lt, ok := tailT.(*ast.ListType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf(":: tail must be a list, got %s", tailT.String()))
		}
		if !typesEqual(headT, lt.Elem) {
			return nil, errAt(e.Pos, fmt.Sprintf(":: head type %s doesn't match list element type %s",
				headT.String(), lt.Elem.String()))
		}
		return tailT, nil

	case *ast.InlExpr:
		t, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		if e.SumType == nil {
			return nil, errAt(e.Pos, "inl requires type annotation: inl expr as LeftType + RightType")
		}
		st, ok := e.SumType.(*ast.SumType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("inl type annotation must be a sum type, got %s", e.SumType.String()))
		}
		if !typesEqual(t, st.Left) {
			return nil, errAt(e.Pos, fmt.Sprintf("inl expression has type %s but sum left type is %s",
				t.String(), st.Left.String()))
		}
		return st, nil

	case *ast.InrExpr:
		t, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		if e.SumType == nil {
			return nil, errAt(e.Pos, "inr requires type annotation: inr expr as LeftType + RightType")
		}
		st, ok := e.SumType.(*ast.SumType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("inr type annotation must be a sum type, got %s", e.SumType.String()))
		}
		if !typesEqual(t, st.Right) {
			return nil, errAt(e.Pos, fmt.Sprintf("inr expression has type %s but sum right type is %s",
				t.String(), st.Right.String()))
		}
		return st, nil

	case *ast.CaseExpr:
		return c.checkCase(e, env)

	case *ast.FixExpr:
		t, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		ft, ok := t.(*ast.FuncType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("fix expects a function type A -> A, got %s", t.String()))
		}
		if !typesEqual(ft.Param, ft.Return) {
			return nil, errAt(e.Pos, fmt.Sprintf("fix expects type A -> A, but got %s -> %s",
				ft.Param.String(), ft.Return.String()))
		}
		return ft.Return, nil

	case *ast.PrintExpr:
		_, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		return &ast.UnitType{}, nil

	case *ast.ImportExpr:
		return &ast.UnitType{}, nil

	case *ast.TupleExpr:
		elems := make([]ast.Type, len(e.Elems))
		for i, elem := range e.Elems {
			t, err := c.check(elem, env)
			if err != nil {
				return nil, err
			}
			elems[i] = t
		}
		return &ast.TupleType{Elems: elems}, nil

	case *ast.TupleAccessExpr:
		t, err := c.check(e.Tuple, env)
		if err != nil {
			return nil, err
		}
		switch tt := t.(type) {
		case *ast.TupleType:
			if e.Index < 0 || e.Index >= len(tt.Elems) {
				return nil, errAt(e.Pos, fmt.Sprintf("tuple index %d out of bounds for tuple of size %d", e.Index, len(tt.Elems)))
			}
			return tt.Elems[e.Index], nil
		case *ast.PairType:
			if e.Index == 0 {
				return tt.First, nil
			} else if e.Index == 1 {
				return tt.Second, nil
			}
			return nil, errAt(e.Pos, fmt.Sprintf("pair index %d out of bounds (pair has 2 elements)", e.Index))
		default:
			return nil, errAt(e.Pos, fmt.Sprintf("index access expects a tuple or pair, got %s", t.String()))
		}

	case *ast.RecordExpr:
		fields := make([]ast.RecordFieldType, len(e.Fields))
		for i, f := range e.Fields {
			t, err := c.check(f.Value, env)
			if err != nil {
				return nil, err
			}
			fields[i] = ast.RecordFieldType{Name: f.Name, Type: t}
		}
		return &ast.RecordType{Fields: fields}, nil

	case *ast.RecordAccessExpr:
		t, err := c.check(e.Record, env)
		if err != nil {
			return nil, err
		}
		rt, ok := t.(*ast.RecordType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("record access expects a record, got %s", t.String()))
		}
		for _, f := range rt.Fields {
			if f.Name == e.Field {
				return f.Type, nil
			}
		}
		return nil, errAt(e.Pos, fmt.Sprintf("record has no field '%s'", e.Field))

	case *ast.ForExpr:
		startT, err := c.check(e.Start, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(startT, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("for loop start must be Int, got %s", startT.String()))
		}
		endT, err := c.check(e.End, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(endT, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("for loop end must be Int, got %s", endT.String()))
		}
		bodyEnv := NewTypeEnv(env)
		bodyEnv.Set(e.Var, &ast.IntType{})
		_, err = c.check(e.Body, bodyEnv)
		if err != nil {
			return nil, err
		}
		return &ast.UnitType{}, nil

	case *ast.LengthExpr:
		t, err := c.check(e.Expr, env)
		if err != nil {
			return nil, err
		}
		switch t.(type) {
		case *ast.StringType, *ast.ListType:
			return &ast.IntType{}, nil
		default:
			return nil, errAt(e.Pos, fmt.Sprintf("length expects a String or List, got %s", t.String()))
		}

	case *ast.TypeDecl:
		if err := c.checkTypeDecl(e); err != nil {
			return nil, err
		}
		return &ast.UnitType{}, nil

	case *ast.CharAtExpr:
		strT, err := c.check(e.Str, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(strT, &ast.StringType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("charAt expects a String, got %s", strT.String()))
		}
		idxT, err := c.check(e.Index, env)
		if err != nil {
			return nil, err
		}
		if !typesEqual(idxT, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("charAt index must be Int, got %s", idxT.String()))
		}
		return &ast.StringType{}, nil

	default:
		return nil, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (c *Checker) checkBinOp(e *ast.BinOp, env *TypeEnv) (ast.Type, error) {
	leftT, err := c.check(e.Left, env)
	if err != nil {
		return nil, err
	}
	rightT, err := c.check(e.Right, env)
	if err != nil {
		return nil, err
	}

	switch e.Op {
	case "+", "-", "*", "/", "%":
		// String concatenation
		if e.Op == "+" {
			if typesEqual(leftT, &ast.StringType{}) && typesEqual(rightT, &ast.StringType{}) {
				return &ast.StringType{}, nil
			}
		}
		if !typesEqual(leftT, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '%s' expects Int, left operand has type %s", e.Op, leftT.String()))
		}
		if !typesEqual(rightT, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '%s' expects Int, right operand has type %s", e.Op, rightT.String()))
		}
		return &ast.IntType{}, nil
	case "==", "!=":
		if !typesEqual(leftT, rightT) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '%s' requires same types, got %s and %s",
				e.Op, leftT.String(), rightT.String()))
		}
		return &ast.BoolType{}, nil
	case "<", ">", "<=", ">=":
		if !typesEqual(leftT, &ast.IntType{}) || !typesEqual(rightT, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '%s' expects Int operands, got %s and %s",
				e.Op, leftT.String(), rightT.String()))
		}
		return &ast.BoolType{}, nil
	case "&&", "||":
		if !typesEqual(leftT, &ast.BoolType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '%s' expects Bool, left operand has type %s", e.Op, leftT.String()))
		}
		if !typesEqual(rightT, &ast.BoolType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '%s' expects Bool, right operand has type %s", e.Op, rightT.String()))
		}
		return &ast.BoolType{}, nil
	}
	return nil, errAt(e.Pos, fmt.Sprintf("unknown operator '%s'", e.Op))
}

func (c *Checker) checkUnaryOp(e *ast.UnaryOp, env *TypeEnv) (ast.Type, error) {
	t, err := c.check(e.Expr, env)
	if err != nil {
		return nil, err
	}
	switch e.Op {
	case "!":
		if !typesEqual(t, &ast.BoolType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("operator '!' expects Bool, got %s", t.String()))
		}
		return &ast.BoolType{}, nil
	case "-":
		if !typesEqual(t, &ast.IntType{}) {
			return nil, errAt(e.Pos, fmt.Sprintf("unary '-' expects Int, got %s", t.String()))
		}
		return &ast.IntType{}, nil
	}
	return nil, errAt(e.Pos, fmt.Sprintf("unknown unary operator '%s'", e.Op))
}

func (c *Checker) checkCase(e *ast.CaseExpr, env *TypeEnv) (ast.Type, error) {
	scrutT, err := c.check(e.Scrutinee, env)
	if err != nil {
		return nil, err
	}

	if len(e.Branches) == 0 {
		return nil, errAt(e.Pos, "case expression must have at least one branch")
	}

	var resultType ast.Type

	// Determine if we're matching on a list or sum
	firstPat := e.Branches[0].Pattern
	switch firstPat.(type) {
	case *ast.NilPattern, *ast.ConsPattern:
		// List case
		lt, ok := scrutT.(*ast.ListType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("case matching on list pattern but scrutinee has type %s", scrutT.String()))
		}
		for _, branch := range e.Branches {
			branchEnv := NewTypeEnv(env)
			switch p := branch.Pattern.(type) {
			case *ast.NilPattern:
				// no bindings
			case *ast.ConsPattern:
				branchEnv.Set(p.Head, lt.Elem)
				branchEnv.Set(p.Tail, scrutT)
			default:
				return nil, errAt(e.Pos, "mixed patterns in case expression")
			}
			bt, err := c.check(branch.Body, branchEnv)
			if err != nil {
				return nil, err
			}
			if resultType == nil {
				resultType = bt
			} else if !typesEqual(resultType, bt) {
				return nil, errAt(e.Pos, fmt.Sprintf("case branches must have same type: expected %s, got %s",
					resultType.String(), bt.String()))
			}
		}
	case *ast.InlPattern, *ast.InrPattern:
		// Sum case
		st, ok := scrutT.(*ast.SumType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("case matching on sum pattern but scrutinee has type %s", scrutT.String()))
		}
		for _, branch := range e.Branches {
			branchEnv := NewTypeEnv(env)
			switch p := branch.Pattern.(type) {
			case *ast.InlPattern:
				branchEnv.Set(p.Name, st.Left)
			case *ast.InrPattern:
				branchEnv.Set(p.Name, st.Right)
			default:
				return nil, errAt(e.Pos, "mixed patterns in case expression")
			}
			bt, err := c.check(branch.Body, branchEnv)
			if err != nil {
				return nil, err
			}
			if resultType == nil {
				resultType = bt
			} else if !typesEqual(resultType, bt) {
				return nil, errAt(e.Pos, fmt.Sprintf("case branches must have same type: expected %s, got %s",
					resultType.String(), bt.String()))
			}
		}
	case *ast.ConstructorPattern:
		// ADT case
		adtT, ok := scrutT.(*ast.ADTType)
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("case matching on constructor pattern but scrutinee has type %s", scrutT.String()))
		}
		adtDef, ok := c.ADTs[adtT.Name]
		if !ok {
			return nil, errAt(e.Pos, fmt.Sprintf("unknown type '%s'", adtT.Name))
		}
		for _, branch := range e.Branches {
			branchEnv := NewTypeEnv(env)
			cp, ok := branch.Pattern.(*ast.ConstructorPattern)
			if !ok {
				return nil, errAt(e.Pos, "mixed patterns in case expression")
			}
			// Find the variant
			var found *ast.VariantDef
			for i := range adtDef.Variants {
				if adtDef.Variants[i].Name == cp.Constructor {
					found = &adtDef.Variants[i]
					break
				}
			}
			if found == nil {
				return nil, errAt(e.Pos, fmt.Sprintf("constructor '%s' is not a variant of type '%s'", cp.Constructor, adtT.Name))
			}
			if found.Payload != nil && cp.Arg == "" {
				return nil, errAt(e.Pos, fmt.Sprintf("constructor '%s' expects an argument", cp.Constructor))
			}
			if found.Payload == nil && cp.Arg != "" {
				return nil, errAt(e.Pos, fmt.Sprintf("constructor '%s' takes no arguments", cp.Constructor))
			}
			if cp.Arg != "" {
				branchEnv.Set(cp.Arg, found.Payload)
			}
			bt, err := c.check(branch.Body, branchEnv)
			if err != nil {
				return nil, err
			}
			if resultType == nil {
				resultType = bt
			} else if !typesEqual(resultType, bt) {
				return nil, errAt(e.Pos, fmt.Sprintf("case branches must have same type: expected %s, got %s",
					resultType.String(), bt.String()))
			}
		}
	default:
		return nil, errAt(e.Pos, "unsupported pattern type in case expression")
	}

	return resultType, nil
}

func (c *Checker) checkTypeDecl(td *ast.TypeDecl) error {
	adtType := &ast.ADTType{Name: td.Name}

	// Register the ADT definition
	c.ADTs[td.Name] = &ADTDef{Name: td.Name, Variants: td.Variants}

	// Register constructors in the type environment
	for _, v := range td.Variants {
		if v.Payload == nil {
			// Nullary constructor: type is just the ADT type
			c.env.Set(v.Name, adtType)
		} else {
			// Unary constructor: type is Payload -> ADTType
			c.env.Set(v.Name, &ast.FuncType{Param: v.Payload, Return: adtType})
		}
	}
	return nil
}
