package evaluator

import (
	"my-programming-language/ast"
	"my-programming-language/checker"
	"my-programming-language/lexer"
	"my-programming-language/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseEvalExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	expr, err := parser.New(toks).ParseSingleExpr()
	require.NoError(t, err)
	return expr
}

func TestEvalArithmeticAndClosure(t *testing.T) {
	expr := mustParseEvalExpr(t, "let x: Int = 10 in (fn (y: Int) => x + y) 5")

	chk := checker.New()
	_, err := chk.Check(expr)
	require.NoError(t, err)

	ev := New()
	val, err := ev.Eval(expr)
	require.NoError(t, err)
	assert.Equal(t, "15", val.String())
}

func TestEvalPrintOutputBuffer(t *testing.T) {
	toks, err := lexer.New(`print "Hello "; println "World!";`).Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)

	chk := checker.New()
	require.NoError(t, chk.CheckProgram(prog))

	ev := New()
	_, err = ev.EvalProgram(prog)
	require.NoError(t, err)
	require.Equal(t, "Hello World!\n", ev.GetOutput())
}

func TestEvalRuntimeErrorDivisionByZero(t *testing.T) {
	expr := mustParseEvalExpr(t, "10 / 0")
	ev := New()
	_, err := ev.Eval(expr)
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "division by zero"), "unexpected runtime error: %v", err)
}

func checkAndEval(t *testing.T, src string) Value {
	t.Helper()
	expr := mustParseEvalExpr(t, src)
	chk := checker.New()
	_, err := chk.Check(expr)
	require.NoError(t, err)
	ev := New()
	val, err := ev.Eval(expr)
	require.NoError(t, err)
	return val
}

func checkAndEvalProgram(t *testing.T, src string) (Value, *Evaluator) {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	prog, err := parser.New(toks).ParseProgram()
	require.NoError(t, err)

	chk := checker.New()
	require.NoError(t, chk.CheckProgram(prog))

	ev := New()
	val, err := ev.EvalProgram(prog)
	require.NoError(t, err)
	return val, ev
}

// --- Tuple tests ---

func TestEvalTupleConstruct(t *testing.T) {
	val := checkAndEval(t, `(1, true, "hello")`)
	tup, ok := val.(*TupleVal)
	require.True(t, ok, "expected TupleVal, got %T", val)
	require.Len(t, tup.Elems, 3)
	assert.Equal(t, "1", tup.Elems[0].String())
	assert.Equal(t, "true", tup.Elems[1].String())
	assert.Equal(t, "hello", tup.Elems[2].String())
}

func TestEvalTupleString(t *testing.T) {
	val := checkAndEval(t, `(1, true, "hi")`)
	assert.Equal(t, `(1, true, hi)`, val.String())
}

func TestEvalTupleAccess(t *testing.T) {
	val := checkAndEval(t, `let t = (10, 20, 30) in t.1`)
	assert.Equal(t, "20", val.String())
}

func TestEvalTupleAccessOnPair(t *testing.T) {
	val := checkAndEval(t, `let p = (42, true) in p.0`)
	assert.Equal(t, "42", val.String())
}

func TestEvalTupleFourElements(t *testing.T) {
	val := checkAndEval(t, `let t = (1, 2, 3, 4) in t.3`)
	assert.Equal(t, "4", val.String())
}

func TestEvalTupleEquality(t *testing.T) {
	val := checkAndEval(t, `(1, 2, 3) == (1, 2, 3)`)
	assert.Equal(t, "true", val.String())

	val2 := checkAndEval(t, `(1, 2, 3) == (1, 2, 4)`)
	assert.Equal(t, "false", val2.String())
}

// --- Record tests ---

func TestEvalRecordConstruct(t *testing.T) {
	val := checkAndEval(t, `{x = 1, y = true}`)
	rec, ok := val.(*RecordVal)
	require.True(t, ok, "expected RecordVal, got %T", val)
	require.Len(t, rec.Fields, 2)
	assert.Equal(t, "x", rec.Fields[0].Name)
	assert.Equal(t, "1", rec.Fields[0].Value.String())
}

func TestEvalRecordString(t *testing.T) {
	val := checkAndEval(t, `{name = "Alice", age = 30}`)
	assert.Equal(t, `{name = Alice, age = 30}`, val.String())
}

func TestEvalRecordAccess(t *testing.T) {
	val := checkAndEval(t, `let r = {x = 42, y = true} in r.x`)
	assert.Equal(t, "42", val.String())
}

func TestEvalRecordAccessSecondField(t *testing.T) {
	val := checkAndEval(t, `let r = {a = 1, b = 2, c = 3} in r.c`)
	assert.Equal(t, "3", val.String())
}

func TestEvalRecordEquality(t *testing.T) {
	val := checkAndEval(t, `{x = 1, y = 2} == {x = 1, y = 2}`)
	assert.Equal(t, "true", val.String())

	val2 := checkAndEval(t, `{x = 1, y = 2} == {x = 1, y = 3}`)
	assert.Equal(t, "false", val2.String())
}

func TestEvalNestedRecord(t *testing.T) {
	val := checkAndEval(t, `let r = {inner = {val = 99}} in r.inner`)
	assert.Equal(t, "{val = 99}", val.String())
}

// --- For loop tests ---

func TestEvalForLoop(t *testing.T) {
	_, ev := checkAndEvalProgram(t, `for i = 0 to 3 do println i end`)
	assert.Equal(t, "0\n1\n2\n", ev.GetOutput())
}

func TestEvalForLoopReturnsUnit(t *testing.T) {
	val := checkAndEval(t, `for i = 0 to 3 do i end`)
	assert.Equal(t, "()", val.String())
}

func TestEvalForLoopEmptyRange(t *testing.T) {
	_, ev := checkAndEvalProgram(t, `for i = 5 to 5 do println i end`)
	assert.Equal(t, "", ev.GetOutput())
}

func TestEvalNestedForLoop(t *testing.T) {
	_, ev := checkAndEvalProgram(t, `for i = 0 to 2 do for j = 0 to 2 do println (i * 10 + j) end end`)
	assert.Equal(t, "0\n1\n10\n11\n", ev.GetOutput())
}

// --- Length tests ---

func TestEvalLengthString(t *testing.T) {
	val := checkAndEval(t, `length "hello"`)
	assert.Equal(t, "5", val.String())
}

func TestEvalLengthEmptyString(t *testing.T) {
	val := checkAndEval(t, `length ""`)
	assert.Equal(t, "0", val.String())
}

func TestEvalLengthList(t *testing.T) {
	val := checkAndEval(t, `length [1, 2, 3]`)
	assert.Equal(t, "3", val.String())
}

func TestEvalLengthEmptyList(t *testing.T) {
	val := checkAndEval(t, `length ([] : [Int])`)
	assert.Equal(t, "0", val.String())
}

// --- CharAt tests ---

func TestEvalCharAt(t *testing.T) {
	val := checkAndEval(t, `charAt "hello" 0`)
	assert.Equal(t, "h", val.String())
}

func TestEvalCharAtLast(t *testing.T) {
	val := checkAndEval(t, `charAt "hello" 4`)
	assert.Equal(t, "o", val.String())
}

func TestEvalCharAtOutOfBounds(t *testing.T) {
	expr := mustParseEvalExpr(t, `charAt "hi" 5`)
	chk := checker.New()
	_, err := chk.Check(expr)
	require.NoError(t, err)
	ev := New()
	_, err = ev.Eval(expr)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "out of bounds")
}

// --- Value String representations ---

func TestTupleValString(t *testing.T) {
	v := &TupleVal{Elems: []Value{&IntVal{1}, &BoolVal{true}}}
	assert.Equal(t, "(1, true)", v.String())
}

// --- ADT tests ---

func TestEvalADTNullaryConstructor(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `type Color = Red | Green | Blue; Red`)
	cv, ok := val.(*ConstructorVal)
	require.True(t, ok, "expected ConstructorVal, got %T", val)
	assert.Equal(t, "Red", cv.Tag)
	assert.Nil(t, cv.Value)
}

func TestEvalADTUnaryConstructor(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `type Option = None | Some of Int; Some 42`)
	cv, ok := val.(*ConstructorVal)
	require.True(t, ok, "expected ConstructorVal, got %T", val)
	assert.Equal(t, "Some", cv.Tag)
	assert.Equal(t, "42", cv.Value.String())
}

func TestEvalADTPatternMatchNullary(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `
type Color = Red | Green | Blue;
let name : Color -> String = fn (c : Color) =>
  case c of
    | Red => "red"
    | Green => "green"
    | Blue => "blue";
name Green
`)
	assert.Equal(t, "green", val.String())
}

func TestEvalADTPatternMatchUnary(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `
type Option = None | Some of Int;
let getValue : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Some n => n;
getValue (Some 99)
`)
	assert.Equal(t, "99", val.String())
}

func TestEvalADTPatternMatchNone(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `
type Option = None | Some of Int;
let getValue : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Some n => n;
getValue None
`)
	assert.Equal(t, "0", val.String())
}

func TestEvalADTMultiplePayloadTypes(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `
type Shape = Circle of Int | Rectangle of (Int, Int);
let area : Shape -> Int = fn (s : Shape) =>
  case s of
    | Circle r => r * r * 3
    | Rectangle dims => fst dims * snd dims;
area (Rectangle (3, 4))
`)
	assert.Equal(t, "12", val.String())
}

func TestEvalADTRecursiveType(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `
type Expr = Lit of Int | Add of (Expr, Expr) | Neg of Expr;
let eval : Expr -> Int = fix fn (eval : Expr -> Int) =>
  fn (e : Expr) =>
    case e of
      | Lit n => n
      | Add pair => eval (fst pair) + eval (snd pair)
      | Neg inner => 0 - eval inner;
eval (Add (Lit 10, Neg (Lit 3)))
`)
	assert.Equal(t, "7", val.String())
}

func TestEvalADTEquality(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `type Color = Red | Green | Blue; Red == Red`)
	assert.Equal(t, "true", val.String())

	val2, _ := checkAndEvalProgram(t, `type Color = Red | Green | Blue; Red == Blue`)
	assert.Equal(t, "false", val2.String())
}

func TestEvalADTEqualityWithPayload(t *testing.T) {
	val, _ := checkAndEvalProgram(t, `type Option = None | Some of Int; (Some 1) == (Some 1)`)
	assert.Equal(t, "true", val.String())

	val2, _ := checkAndEvalProgram(t, `type Option = None | Some of Int; (Some 1) == (Some 2)`)
	assert.Equal(t, "false", val2.String())

	val3, _ := checkAndEvalProgram(t, `type Option = None | Some of Int; None == None`)
	assert.Equal(t, "true", val3.String())
}

func TestEvalADTPrint(t *testing.T) {
	_, ev := checkAndEvalProgram(t, `
type Option = None | Some of Int;
println None;
println (Some 42);
`)
	assert.Equal(t, "None\nSome 42\n", ev.GetOutput())
}

func TestConstructorValString(t *testing.T) {
	v := &ConstructorVal{Tag: "None", Value: nil}
	assert.Equal(t, "None", v.String())

	v2 := &ConstructorVal{Tag: "Some", Value: &IntVal{Value: 42}}
	assert.Equal(t, "Some 42", v2.String())
}

func TestRecordValString(t *testing.T) {
	v := &RecordVal{Fields: []RecordFieldVal{
		{Name: "x", Value: &IntVal{42}},
	}}
	assert.Equal(t, "{x = 42}", v.String())
}
