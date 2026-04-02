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

func TestTokenizeNewKeywords(t *testing.T) {
	src := `for i = 0 to 10 do println i end`
	toks, err := New(src).Tokenize()
	require.NoError(t, err)

	wantTypes := []token.TokenType{
		token.FOR, token.IDENT, token.EQ, token.INT, token.TO, token.INT,
		token.DO, token.PRINTLN, token.IDENT, token.END, token.EOF,
	}
	require.Len(t, toks, len(wantTypes))
	for i, want := range wantTypes {
		assert.Equal(t, want, toks[i].Type, "token %d: expected %s, got %s", i, want, toks[i].Type)
	}
}

func TestTokenizeLengthAndCharAt(t *testing.T) {
	src := `length "hi"; charAt "hi" 0`
	toks, err := New(src).Tokenize()
	require.NoError(t, err)

	wantTypes := []token.TokenType{
		token.LENGTH, token.STRING, token.SEMICOLON,
		token.CHARAT, token.STRING, token.INT, token.EOF,
	}
	require.Len(t, toks, len(wantTypes))
	for i, want := range wantTypes {
		assert.Equal(t, want, toks[i].Type, "token %d", i)
	}
}

func TestTokenizeTypeDecl(t *testing.T) {
	src := `type Option = None | Some of Int`
	toks, err := New(src).Tokenize()
	require.NoError(t, err)

	wantTypes := []token.TokenType{
		token.TYPE, token.IDENT, token.EQ, token.IDENT, token.PIPE,
		token.IDENT, token.OF, token.INT_TYPE, token.EOF,
	}
	require.Len(t, toks, len(wantTypes))
	for i, want := range wantTypes {
		assert.Equal(t, want, toks[i].Type, "token %d: expected %s, got %s", i, want, toks[i].Type)
	}
}

func TestTokenizeRecordAndDot(t *testing.T) {
	src := `{x = 1}.x`
	toks, err := New(src).Tokenize()
	require.NoError(t, err)

	wantTypes := []token.TokenType{
		token.LBRACE, token.IDENT, token.EQ, token.INT, token.RBRACE,
		token.DOT, token.IDENT, token.EOF,
	}
	require.Len(t, toks, len(wantTypes))
	for i, want := range wantTypes {
		assert.Equal(t, want, toks[i].Type, "token %d", i)
	}
}
