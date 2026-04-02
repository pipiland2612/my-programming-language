package lexer

import (
	"my-programming-language/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizeKeywordsOperatorsAndComments(t *testing.T) {
	src := `
// comment
let x: Int = 1 + 2;
println "ok";
/* block */
`
	toks, err := New(src).Tokenize()
	require.NoError(t, err)

	wantTypes := []token.TokenType{
		token.LET, token.IDENT, token.COLON, token.INT_TYPE, token.EQ, token.INT, token.PLUS, token.INT, token.SEMICOLON,
		token.PRINTLN, token.STRING, token.SEMICOLON, token.EOF,
	}
	require.Len(t, toks, len(wantTypes))
	for i, want := range wantTypes {
		require.Equal(t, want, toks[i].Type, "token %d type mismatch", i)
	}
}

func TestTokenizeUnterminatedString(t *testing.T) {
	_, err := New(`"abc`).Tokenize()
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "unterminated string"), "unexpected error: %v", err)
}
