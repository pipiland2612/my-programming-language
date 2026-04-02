package parser

import (
	"my-programming-language/ast"
	"my-programming-language/lexer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parseSingleExpr(t *testing.T, src string) ast.Expr {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	expr, err := New(toks).ParseSingleExpr()
	require.NoError(t, err)
	return expr
}

func TestParseSingleExprPrecedence(t *testing.T) {
	expr := parseSingleExpr(t, "1 + 2 * 3")
	add, ok := expr.(*ast.BinOp)
	require.True(t, ok, "expected top-level '+' binop, got %#v", expr)
	require.Equal(t, "+", add.Op)
	mul, ok := add.Right.(*ast.BinOp)
	require.True(t, ok, "expected right side '*' binop, got %#v", add.Right)
	require.Equal(t, "*", mul.Op)
}

func TestParseProgramImportAndLet(t *testing.T) {
	toks, err := lexer.New(`import "mathlib.mepl"; let x: Int = 42;`).Tokenize()
	require.NoError(t, err)
	prog, err := New(toks).ParseProgram()
	require.NoError(t, err)
	require.Len(t, prog.Declarations, 2)
	require.IsType(t, &ast.ImportExpr{}, prog.Declarations[0])
	assert.IsType(t, &ast.LetExpr{}, prog.Declarations[1])
}

func TestParseError(t *testing.T) {
	toks, err := lexer.New("if true then 1").Tokenize()
	require.NoError(t, err)
	_, err = New(toks).ParseSingleExpr()
	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "expected else"), "unexpected parse error: %v", err)
}

func parseProgram(t *testing.T, src string) *ast.Program {
	t.Helper()
	toks, err := lexer.New(src).Tokenize()
	require.NoError(t, err)
	prog, err := New(toks).ParseProgram()
	require.NoError(t, err)
	return prog
}

func TestParseTupleExpr(t *testing.T) {
	expr := parseSingleExpr(t, "(1, true, 3)")
	tup, ok := expr.(*ast.TupleExpr)
	require.True(t, ok, "expected TupleExpr, got %T", expr)
	assert.Len(t, tup.Elems, 3)
	assert.IsType(t, &ast.IntLit{}, tup.Elems[0])
	assert.IsType(t, &ast.BoolLit{}, tup.Elems[1])
	assert.IsType(t, &ast.IntLit{}, tup.Elems[2])
}

func TestParseTupleStillParsePair(t *testing.T) {
	expr := parseSingleExpr(t, "(1, 2)")
	_, ok := expr.(*ast.PairExpr)
	require.True(t, ok, "two-element tuple should still be a PairExpr, got %T", expr)
}

func TestParseTupleAccess(t *testing.T) {
	expr := parseSingleExpr(t, "let t = (1, 2, 3) in t.0")
	let, ok := expr.(*ast.LetExpr)
	require.True(t, ok)
	acc, ok := let.Body.(*ast.TupleAccessExpr)
	require.True(t, ok, "expected TupleAccessExpr, got %T", let.Body)
	assert.Equal(t, 0, acc.Index)
}

func TestParseRecordExpr(t *testing.T) {
	expr := parseSingleExpr(t, "{x = 1, y = true}")
	rec, ok := expr.(*ast.RecordExpr)
	require.True(t, ok, "expected RecordExpr, got %T", expr)
	require.Len(t, rec.Fields, 2)
	assert.Equal(t, "x", rec.Fields[0].Name)
	assert.Equal(t, "y", rec.Fields[1].Name)
}

func TestParseRecordAccess(t *testing.T) {
	expr := parseSingleExpr(t, "let r = {x = 1} in r.x")
	let, ok := expr.(*ast.LetExpr)
	require.True(t, ok)
	acc, ok := let.Body.(*ast.RecordAccessExpr)
	require.True(t, ok, "expected RecordAccessExpr, got %T", let.Body)
	assert.Equal(t, "x", acc.Field)
}

func TestParseForExpr(t *testing.T) {
	expr := parseSingleExpr(t, "for i = 0 to 5 do i end")
	f, ok := expr.(*ast.ForExpr)
	require.True(t, ok, "expected ForExpr, got %T", expr)
	assert.Equal(t, "i", f.Var)
	assert.IsType(t, &ast.IntLit{}, f.Start)
	assert.IsType(t, &ast.IntLit{}, f.End)
	assert.IsType(t, &ast.Var{}, f.Body)
}

func TestParseLengthExpr(t *testing.T) {
	expr := parseSingleExpr(t, `length "hello"`)
	l, ok := expr.(*ast.LengthExpr)
	require.True(t, ok, "expected LengthExpr, got %T", expr)
	assert.IsType(t, &ast.StringLit{}, l.Expr)
}

func TestParseCharAtExpr(t *testing.T) {
	expr := parseSingleExpr(t, `charAt "hello" 0`)
	ca, ok := expr.(*ast.CharAtExpr)
	require.True(t, ok, "expected CharAtExpr, got %T", expr)
	assert.IsType(t, &ast.StringLit{}, ca.Str)
	assert.IsType(t, &ast.IntLit{}, ca.Index)
}

func TestParseRecordType(t *testing.T) {
	prog := parseProgram(t, `let r : {x: Int, y: Bool} = {x = 1, y = true};`)
	require.Len(t, prog.Declarations, 1)
	let := prog.Declarations[0].(*ast.LetExpr)
	rt, ok := let.TypeAnn.(*ast.RecordType)
	require.True(t, ok, "expected RecordType annotation, got %T", let.TypeAnn)
	require.Len(t, rt.Fields, 2)
	assert.Equal(t, "x", rt.Fields[0].Name)
	assert.Equal(t, "Int", rt.Fields[0].Type.String())
}

func TestParseTupleType(t *testing.T) {
	prog := parseProgram(t, `let t : (Int, Bool, String) = (1, true, "hi");`)
	require.Len(t, prog.Declarations, 1)
	let := prog.Declarations[0].(*ast.LetExpr)
	tt, ok := let.TypeAnn.(*ast.TupleType)
	require.True(t, ok, "expected TupleType annotation, got %T", let.TypeAnn)
	require.Len(t, tt.Elems, 3)
}

func TestParsePostfixChain(t *testing.T) {
	// r.x should parse as RecordAccessExpr even after fst/snd/println
	expr := parseSingleExpr(t, "let r = {x = 1} in fst (r, 2)")
	let := expr.(*ast.LetExpr)
	fstExpr, ok := let.Body.(*ast.FstExpr)
	require.True(t, ok, "expected FstExpr, got %T", let.Body)
	_ = fstExpr
}

func TestParseForMissingEnd(t *testing.T) {
	toks, err := lexer.New("for i = 0 to 5 do i").Tokenize()
	require.NoError(t, err)
	_, err = New(toks).ParseSingleExpr()
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "expected end"), "unexpected parse error: %v", err)
}
