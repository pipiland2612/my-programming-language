package token

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenTypeString(t *testing.T) {
	require.Equal(t, "+", PLUS.String(), "expected PLUS to stringify to '+'")

	unknown := TokenType(9999)
	require.Equal(t, "UNKNOWN", unknown.String(), "expected unknown token to stringify to UNKNOWN")
}

func TestKeywordsMap(t *testing.T) {
	cases := map[string]TokenType{
		"let":     LET,
		"println": PRINTLN,
		"Int":     INT_TYPE,
		"Bool":    BOOL_TYPE,
		"for":     FOR,
		"to":      TO,
		"do":      DO,
		"end":     END,
		"length":  LENGTH,
		"charAt":  CHARAT,
	}

	for kw, want := range cases {
		got, ok := Keywords[kw]
		require.True(t, ok, "keyword %q not found", kw)
		assert.Equal(t, want, got, "keyword %q mismatch", kw)
	}
}

func TestNewKeywordTokenStrings(t *testing.T) {
	assert.Equal(t, "for", FOR.String())
	assert.Equal(t, "to", TO.String())
	assert.Equal(t, "do", DO.String())
	assert.Equal(t, "end", END.String())
	assert.Equal(t, "length", LENGTH.String())
	assert.Equal(t, "charAt", CHARAT.String())
}

func TestTokenStructFields(t *testing.T) {
	tok := Token{
		Type:    IDENT,
		Literal: "abc",
		Pos: Pos{
			Line:   3,
			Column: 7,
		},
	}

	require.Equal(t, IDENT, tok.Type)
	require.Equal(t, "abc", tok.Literal)
	require.Equal(t, 3, tok.Pos.Line)
	require.Equal(t, 7, tok.Pos.Column)
}
