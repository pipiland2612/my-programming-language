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

func TestTupleTypeString(t *testing.T) {
	tt := TupleType{Elems: []Type{IntType{}, BoolType{}, StringType{}}}
	assert.Equal(t, "(Int, Bool, String)", tt.String())
}

func TestRecordTypeString(t *testing.T) {
	rt := RecordType{Fields: []RecordFieldType{
		{Name: "x", Type: IntType{}},
		{Name: "y", Type: BoolType{}},
	}}
	assert.Equal(t, "{x: Int, y: Bool}", rt.String())
}

func TestNewExprGetPos(t *testing.T) {
	pos := token.Pos{Line: 1, Column: 1}

	assert.Equal(t, pos, TupleExpr{Pos: pos}.GetPos())
	assert.Equal(t, pos, TupleAccessExpr{Pos: pos}.GetPos())
	assert.Equal(t, pos, RecordExpr{Pos: pos}.GetPos())
	assert.Equal(t, pos, RecordAccessExpr{Pos: pos}.GetPos())
	assert.Equal(t, pos, ForExpr{Pos: pos}.GetPos())
	assert.Equal(t, pos, LengthExpr{Pos: pos}.GetPos())
	assert.Equal(t, pos, CharAtExpr{Pos: pos}.GetPos())
}
