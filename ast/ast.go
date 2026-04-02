package ast

import (
	"my-programming-language/token"
	"strings"
)

// Types

type Type interface {
	typeNode()
	String() string
}

type IntType struct{}
type BoolType struct{}
type StringType struct{}
type UnitType struct{}
type FuncType struct{ Param, Return Type }
type PairType struct{ First, Second Type }
type ListType struct{ Elem Type }
type SumType struct{ Left, Right Type }
type TupleType struct{ Elems []Type }
type RecordType struct{ Fields []RecordFieldType }
type RecordFieldType struct {
	Name string
	Type Type
}
type ADTType struct{ Name string }

func (IntType) typeNode()    {}
func (BoolType) typeNode()   {}
func (StringType) typeNode() {}
func (UnitType) typeNode()   {}
func (FuncType) typeNode()   {}
func (PairType) typeNode()   {}
func (ListType) typeNode()   {}
func (SumType) typeNode()    {}
func (TupleType) typeNode()  {}
func (RecordType) typeNode() {}
func (ADTType) typeNode()    {}

func (IntType) String() string    { return "Int" }
func (BoolType) String() string   { return "Bool" }
func (StringType) String() string { return "String" }
func (UnitType) String() string   { return "Unit" }
func (t FuncType) String() string {
	p := t.Param.String()
	if _, ok := t.Param.(FuncType); ok {
		p = "(" + p + ")"
	}
	return p + " -> " + t.Return.String()
}
func (t PairType) String() string { return "(" + t.First.String() + ", " + t.Second.String() + ")" }
func (t ListType) String() string { return "[" + t.Elem.String() + "]" }
func (t SumType) String() string { return t.Left.String() + " + " + t.Right.String() }
func (t TupleType) String() string {
	parts := make([]string, len(t.Elems))
	for i, e := range t.Elems {
		parts[i] = e.String()
	}
	return "(" + strings.Join(parts, ", ") + ")"
}
func (t RecordType) String() string {
	parts := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		parts[i] = f.Name + ": " + f.Type.String()
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
func (t ADTType) String() string { return t.Name }

// Expressions

type Expr interface {
	exprNode()
	GetPos() token.Pos
}

type IntLit struct {
	Value int
	Pos   token.Pos
}

type BoolLit struct {
	Value bool
	Pos   token.Pos
}

type StringLit struct {
	Value string
	Pos   token.Pos
}

type UnitLit struct {
	Pos token.Pos
}

type Var struct {
	Name string
	Pos  token.Pos
}

type BinOp struct {
	Op    string
	Left  Expr
	Right Expr
	Pos   token.Pos
}

type UnaryOp struct {
	Op   string
	Expr Expr
	Pos  token.Pos
}

type IfExpr struct {
	Cond Expr
	Then Expr
	Else Expr
	Pos  token.Pos
}

type LetExpr struct {
	Name    string
	TypeAnn Type // optional
	Value   Expr
	Body    Expr // nil for top-level declarations
	Pos     token.Pos
}

type FnExpr struct {
	Param     string
	ParamType Type
	Body      Expr
	Pos       token.Pos
}

type AppExpr struct {
	Func Expr
	Arg  Expr
	Pos  token.Pos
}

type PairExpr struct {
	First  Expr
	Second Expr
	Pos    token.Pos
}

type FstExpr struct {
	Expr Expr
	Pos  token.Pos
}

type SndExpr struct {
	Expr Expr
	Pos  token.Pos
}

type ListExpr struct {
	Elems    []Expr
	ElemType Type // optional type annotation for empty lists
	Pos      token.Pos
}

type ConsExpr struct {
	Head Expr
	Tail Expr
	Pos  token.Pos
}

type InlExpr struct {
	Expr    Expr
	SumType Type // type annotation: the full sum type
	Pos     token.Pos
}

type InrExpr struct {
	Expr    Expr
	SumType Type // type annotation: the full sum type
	Pos     token.Pos
}

type CaseExpr struct {
	Scrutinee Expr
	Branches  []CaseBranch
	Pos       token.Pos
}

type CaseBranch struct {
	Pattern Pattern
	Body    Expr
}

type Pattern interface {
	patternNode()
}

type NilPattern struct{}                          // []
type ConsPattern struct{ Head, Tail string }      // h :: t
type InlPattern struct{ Name string }             // inl x
type InrPattern struct{ Name string }             // inr y
type ConstructorPattern struct {
	Constructor string // e.g. "Some", "None"
	Arg         string // bound variable name, empty for nullary
}

func (NilPattern) patternNode()          {}
func (ConsPattern) patternNode()         {}
func (InlPattern) patternNode()          {}
func (InrPattern) patternNode()          {}
func (ConstructorPattern) patternNode()  {}

type FixExpr struct {
	Expr Expr
	Pos  token.Pos
}

type PrintExpr struct {
	Expr    Expr
	Newline bool
	Pos     token.Pos
}

type ImportExpr struct {
	Path string
	Pos  token.Pos
}

type TupleExpr struct {
	Elems []Expr
	Pos   token.Pos
}

type TupleAccessExpr struct {
	Tuple Expr
	Index int
	Pos   token.Pos
}

type RecordExpr struct {
	Fields []RecordField
	Pos    token.Pos
}

type RecordField struct {
	Name  string
	Value Expr
}

type RecordAccessExpr struct {
	Record Expr
	Field  string
	Pos    token.Pos
}

type ForExpr struct {
	Var   string
	Start Expr
	End   Expr
	Body  Expr
	Pos   token.Pos
}

type TypeDecl struct {
	Name     string
	Variants []VariantDef
	Pos      token.Pos
}

type VariantDef struct {
	Name    string
	Payload Type // nil for nullary constructors
}

type LengthExpr struct {
	Expr Expr
	Pos  token.Pos
}

type CharAtExpr struct {
	Str   Expr
	Index Expr
	Pos   token.Pos
}

// exprNode implementations
func (IntLit) exprNode()    {}
func (BoolLit) exprNode()   {}
func (StringLit) exprNode() {}
func (UnitLit) exprNode()   {}
func (Var) exprNode()       {}
func (BinOp) exprNode()     {}
func (UnaryOp) exprNode()   {}
func (IfExpr) exprNode()    {}
func (LetExpr) exprNode()   {}
func (FnExpr) exprNode()    {}
func (AppExpr) exprNode()   {}
func (PairExpr) exprNode()  {}
func (FstExpr) exprNode()   {}
func (SndExpr) exprNode()   {}
func (ListExpr) exprNode()  {}
func (ConsExpr) exprNode()  {}
func (InlExpr) exprNode()   {}
func (InrExpr) exprNode()   {}
func (CaseExpr) exprNode()  {}
func (FixExpr) exprNode()   {}
func (PrintExpr) exprNode()        {}
func (ImportExpr) exprNode()       {}
func (TupleExpr) exprNode()        {}
func (TupleAccessExpr) exprNode()  {}
func (RecordExpr) exprNode()       {}
func (RecordAccessExpr) exprNode() {}
func (ForExpr) exprNode()          {}
func (TypeDecl) exprNode()         {}
func (LengthExpr) exprNode()       {}
func (CharAtExpr) exprNode()       {}

// GetPos implementations
func (e IntLit) GetPos() token.Pos    { return e.Pos }
func (e BoolLit) GetPos() token.Pos   { return e.Pos }
func (e StringLit) GetPos() token.Pos { return e.Pos }
func (e UnitLit) GetPos() token.Pos   { return e.Pos }
func (e Var) GetPos() token.Pos       { return e.Pos }
func (e BinOp) GetPos() token.Pos     { return e.Pos }
func (e UnaryOp) GetPos() token.Pos   { return e.Pos }
func (e IfExpr) GetPos() token.Pos    { return e.Pos }
func (e LetExpr) GetPos() token.Pos   { return e.Pos }
func (e FnExpr) GetPos() token.Pos    { return e.Pos }
func (e AppExpr) GetPos() token.Pos   { return e.Pos }
func (e PairExpr) GetPos() token.Pos  { return e.Pos }
func (e FstExpr) GetPos() token.Pos   { return e.Pos }
func (e SndExpr) GetPos() token.Pos   { return e.Pos }
func (e ListExpr) GetPos() token.Pos  { return e.Pos }
func (e ConsExpr) GetPos() token.Pos  { return e.Pos }
func (e InlExpr) GetPos() token.Pos   { return e.Pos }
func (e InrExpr) GetPos() token.Pos   { return e.Pos }
func (e CaseExpr) GetPos() token.Pos  { return e.Pos }
func (e FixExpr) GetPos() token.Pos   { return e.Pos }
func (e PrintExpr) GetPos() token.Pos        { return e.Pos }
func (e ImportExpr) GetPos() token.Pos       { return e.Pos }
func (e TupleExpr) GetPos() token.Pos        { return e.Pos }
func (e TupleAccessExpr) GetPos() token.Pos  { return e.Pos }
func (e RecordExpr) GetPos() token.Pos       { return e.Pos }
func (e RecordAccessExpr) GetPos() token.Pos { return e.Pos }
func (e ForExpr) GetPos() token.Pos          { return e.Pos }
func (e TypeDecl) GetPos() token.Pos         { return e.Pos }
func (e LengthExpr) GetPos() token.Pos       { return e.Pos }
func (e CharAtExpr) GetPos() token.Pos       { return e.Pos }

// Program is a list of top-level declarations
type Program struct {
	Declarations []Expr // LetExpr with Body == nil, or ImportExpr
}
