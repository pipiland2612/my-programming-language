package ast

import (
	"my-programming-language/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeStringFormatting(t *testing.T) {
	fn := FuncType{
		Param:  FuncType{Param: IntType{}, Return: BoolType{}},
		Return: StringType{},
	}
	assert.Equal(t, "(Int -> Bool) -> String", fn.String())

	pair := PairType{First: IntType{}, Second: BoolType{}}
	require.Equal(t, "(Int, Bool)", pair.String())
}

func TestExprGetPos(t *testing.T) {
	pos := token.Pos{Line: 9, Column: 4}
	expr := StringLit{Value: "x", Pos: pos}
	require.Equal(t, pos, expr.GetPos())
}
