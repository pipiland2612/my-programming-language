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
