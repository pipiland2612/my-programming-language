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
