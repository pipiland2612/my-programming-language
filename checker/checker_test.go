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
