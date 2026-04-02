package checker

import (
	"my-programming-language/ast"
	"my-programming-language/lexer"
	"my-programming-language/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	expr, err := parser.New(toks).ParseSingleExpr()
	require.NoError(t, err)
	return expr
}

func TestCheckValidExpression(t *testing.T) {
	expr := mustParseExpr(t, "let f: Int -> Int = fn (x: Int) => x in f 10")
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Int", typ.String())
}

func TestCheckTypeError(t *testing.T) {
	expr := mustParseExpr(t, "1 + true")
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "operator '+' expects Int"), "unexpected type error: %v", err)
}

func TestCheckProgramPopulatesEnvironment(t *testing.T) {
	toks, err := lexer.New("let x: Int = 5; let y: Int = x + 1;").Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)

	chk := New()
	require.NoError(t, chk.CheckProgram(prog))

	tpe, ok := chk.Env().Get("y")
	require.True(t, ok, "expected y in environment")
	require.Equal(t, "Int", tpe.String())
}

// --- Tuple tests ---

func TestCheckTupleExpr(t *testing.T) {
	expr := mustParseExpr(t, `(1, true, "hi")`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, `(Int, Bool, String)`, typ.String())
}

func TestCheckTupleAccess(t *testing.T) {
	expr := mustParseExpr(t, `let t = (1, true, "hi") in t.2`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "String", typ.String())
}

func TestCheckTupleAccessOutOfBounds(t *testing.T) {
	expr := mustParseExpr(t, `let t = (1, true, "hi") in t.5`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "out of bounds")
}

func TestCheckTupleAccessOnPair(t *testing.T) {
	expr := mustParseExpr(t, `let p = (1, true) in p.0`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Int", typ.String())
}

func TestCheckTupleAccessOnNonTuple(t *testing.T) {
	expr := mustParseExpr(t, `let x = 42 in x.0`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expects a tuple or pair")
}

// --- Record tests ---

func TestCheckRecordExpr(t *testing.T) {
	expr := mustParseExpr(t, `{x = 1, y = true}`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "{x: Int, y: Bool}", typ.String())
}

func TestCheckRecordAccess(t *testing.T) {
	expr := mustParseExpr(t, `let r = {name = "Alice", age = 30} in r.age`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Int", typ.String())
}

func TestCheckRecordAccessMissingField(t *testing.T) {
	expr := mustParseExpr(t, `let r = {x = 1} in r.y`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no field 'y'")
}

func TestCheckRecordAccessOnNonRecord(t *testing.T) {
	expr := mustParseExpr(t, `let x = 42 in x.foo`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expects a record")
}

// --- For loop tests ---

func TestCheckForExpr(t *testing.T) {
	expr := mustParseExpr(t, `for i = 0 to 10 do i end`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Unit", typ.String())
}

func TestCheckForExprNonIntStart(t *testing.T) {
	expr := mustParseExpr(t, `for i = true to 10 do i end`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "start must be Int")
}

func TestCheckForExprNonIntEnd(t *testing.T) {
	expr := mustParseExpr(t, `for i = 0 to "hi" do i end`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "end must be Int")
}

func TestCheckForExprBodyUsesLoopVar(t *testing.T) {
	// The loop variable should be typed as Int inside the body
	expr := mustParseExpr(t, `for i = 0 to 5 do i + 1 end`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Unit", typ.String())
}

// --- Length tests ---

func TestCheckLengthString(t *testing.T) {
	expr := mustParseExpr(t, `length "hello"`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Int", typ.String())
}

func TestCheckLengthList(t *testing.T) {
	expr := mustParseExpr(t, `length [1, 2, 3]`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "Int", typ.String())
}

func TestCheckLengthOnInt(t *testing.T) {
	expr := mustParseExpr(t, `length 42`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expects a String or List")
}

// --- CharAt tests ---

func TestCheckCharAt(t *testing.T) {
	expr := mustParseExpr(t, `charAt "hello" 0`)
	chk := New()
	typ, err := chk.Check(expr)
	require.NoError(t, err)
	assert.Equal(t, "String", typ.String())
}

func TestCheckCharAtNonString(t *testing.T) {
	expr := mustParseExpr(t, `charAt 42 0`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expects a String")
}

func TestCheckCharAtNonIntIndex(t *testing.T) {
	expr := mustParseExpr(t, `charAt "hello" true`)
	chk := New()
	_, err := chk.Check(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "index must be Int")
}

// --- typesEqual for new types ---

func TestCheckTupleTypeEquality(t *testing.T) {
	// Two tuples with same structure should match in type annotation
	toks, err := lexer.New(`let t : (Int, Bool, String) = (1, true, "x");`).Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)

	chk := New()
	require.NoError(t, chk.CheckProgram(prog))
}

// --- ADT tests ---

func mustCheckProgram(t *testing.T, src string) *Checker {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)
	chk := New()
	require.NoError(t, chk.CheckProgram(prog))
	return chk
}

func mustFailCheckProgram(t *testing.T, src string) error {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)
	chk := New()
	err = chk.CheckProgram(prog)
	require.Error(t, err)
	return err
}

func TestCheckADTNullaryConstructor(t *testing.T) {
	chk := mustCheckProgram(t, `type Color = Red | Green | Blue; let c : Color = Red;`)
	typ, ok := chk.Env().Get("c")
	require.True(t, ok)
	assert.Equal(t, "Color", typ.String())
}

func TestCheckADTUnaryConstructor(t *testing.T) {
	chk := mustCheckProgram(t, `type Option = None | Some of Int; let x : Option = Some 42;`)
	typ, ok := chk.Env().Get("x")
	require.True(t, ok)
	assert.Equal(t, "Option", typ.String())
}

func TestCheckADTConstructorTypeError(t *testing.T) {
	err := mustFailCheckProgram(t, `type Option = None | Some of Int; let x : Option = Some true;`)
	assert.Contains(t, err.Error(), "type error")
}

func TestCheckADTPatternMatch(t *testing.T) {
	chk := mustCheckProgram(t, `
type Option = None | Some of Int;
let getValue : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Some n => n;
`)
	typ, ok := chk.Env().Get("getValue")
	require.True(t, ok)
	assert.Equal(t, "Option -> Int", typ.String())
}

func TestCheckADTPatternMatchBranchTypeMismatch(t *testing.T) {
	err := mustFailCheckProgram(t, `
type Option = None | Some of Int;
let getValue : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Some n => true;
`)
	assert.Contains(t, err.Error(), "case branches must have same type")
}

func TestCheckADTUnknownConstructorInPattern(t *testing.T) {
	err := mustFailCheckProgram(t, `
type Option = None | Some of Int;
let f : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Foo n => n;
`)
	assert.Contains(t, err.Error(), "not a variant")
}

func TestCheckADTMissingArgInPattern(t *testing.T) {
	err := mustFailCheckProgram(t, `
type Option = None | Some of Int;
let f : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Some => 1;
`)
	assert.Contains(t, err.Error(), "expects an argument")
}

func TestCheckADTExtraArgInPattern(t *testing.T) {
	err := mustFailCheckProgram(t, `
type Color = Red | Green | Blue;
let f : Color -> Int = fn (c : Color) =>
  case c of
    | Red x => 0
    | Green => 1
    | Blue => 2;
`)
	assert.Contains(t, err.Error(), "takes no arguments")
}

func TestCheckADTRecursiveType(t *testing.T) {
	mustCheckProgram(t, `
type Expr = Lit of Int | Neg of Expr;
let e : Expr = Neg (Lit 5);
`)
}

func TestCheckADTTypeEquality(t *testing.T) {
	assert.True(t, typesEqual(&ast.ADTType{Name: "Foo"}, &ast.ADTType{Name: "Foo"}))
	assert.False(t, typesEqual(&ast.ADTType{Name: "Foo"}, &ast.ADTType{Name: "Bar"}))
	assert.False(t, typesEqual(&ast.ADTType{Name: "Foo"}, &ast.IntType{}))
}

func TestCheckRecordTypeEquality(t *testing.T) {
	toks, err := lexer.New(`let r : {x: Int, y: Bool} = {x = 1, y = true};`).Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)

	chk := New()
	require.NoError(t, chk.CheckProgram(prog))
}
