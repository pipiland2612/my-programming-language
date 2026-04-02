package parser

import (
	"fmt"
	"my-programming-language/ast"
	"my-programming-language/token"
	"strconv"
)

type Parser struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) cur() token.Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return token.Token{Type: token.EOF}
}

func (p *Parser) peek() token.Token {
	if p.pos+1 < len(p.tokens) {
		return p.tokens[p.pos+1]
	}
	return token.Token{Type: token.EOF}
}

func (p *Parser) advance() token.Token {
	t := p.cur()
	p.pos++
	return t
}

func (p *Parser) expect(typ token.TokenType) (token.Token, error) {
	t := p.cur()
	if t.Type != typ {
		return t, fmt.Errorf("%d:%d: expected %s, got '%s'", t.Pos.Line, t.Pos.Column, typ.String(), t.Literal)
	}
	p.advance()
	return t, nil
}

func (p *Parser) errAt(pos token.Pos, msg string) error {
	return fmt.Errorf("%d:%d: %s", pos.Line, pos.Column, msg)
}

// ParseProgram parses top-level declarations
func (p *Parser) ParseProgram() (*ast.Program, error) {
	prog := &ast.Program{}
	for p.cur().Type != token.EOF {
		// Skip semicolons between declarations
		for p.cur().Type == token.SEMICOLON {
			p.advance()
		}
		if p.cur().Type == token.EOF {
			break
		}

		if p.cur().Type == token.IMPORT {
			imp, err := p.parseImport()
			if err != nil {
				return nil, err
			}
			prog.Declarations = append(prog.Declarations, imp)
		} else if p.cur().Type == token.LET {
			decl, err := p.parseLetDeclaration()
			if err != nil {
				return nil, err
			}
			prog.Declarations = append(prog.Declarations, decl)
		} else {
			// Allow bare expressions at top level (e.g. println, function calls)
			expr, err := p.ParseExpr()
			if err != nil {
				return nil, err
			}
			prog.Declarations = append(prog.Declarations, expr)
		}

		// Optional semicolon after declaration
		if p.cur().Type == token.SEMICOLON {
			p.advance()
		}
	}
	return prog, nil
}

func (p *Parser) parseImport() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // import
	t, err := p.expect(token.STRING)
	if err != nil {
		return nil, err
	}
	return &ast.ImportExpr{Path: t.Literal, Pos: pos}, nil
}

func (p *Parser) parseLetDeclaration() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // let

	name, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}

	var typeAnn ast.Type
	if p.cur().Type == token.COLON {
		p.advance()
		typeAnn, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	if _, err := p.expect(token.EQ); err != nil {
		return nil, err
	}

	value, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	// Check if this is a let-in expression or a top-level declaration
	if p.cur().Type == token.IN {
		p.advance()
		body, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.LetExpr{Name: name.Literal, TypeAnn: typeAnn, Value: value, Body: body, Pos: pos}, nil
	}

	// Top-level declaration (body is nil)
	return &ast.LetExpr{Name: name.Literal, TypeAnn: typeAnn, Value: value, Body: nil, Pos: pos}, nil
}

// ParseExpr is the entry point for expression parsing
func (p *Parser) ParseExpr() (ast.Expr, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (ast.Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.OR {
		pos := p.cur().Pos
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: "||", Left: left, Right: right, Pos: pos}
	}
	return left, nil
}

func (p *Parser) parseAnd() (ast.Expr, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.AND {
		pos := p.cur().Pos
		p.advance()
		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: "&&", Left: left, Right: right, Pos: pos}
	}
	return left, nil
}

func (p *Parser) parseEquality() (ast.Expr, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.EQEQ || p.cur().Type == token.NEQ {
		pos := p.cur().Pos
		op := p.advance().Literal
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right, Pos: pos}
	}
	return left, nil
}

func (p *Parser) parseComparison() (ast.Expr, error) {
	left, err := p.parseCons()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.LT || p.cur().Type == token.GT ||
		p.cur().Type == token.LEQ || p.cur().Type == token.GEQ {
		pos := p.cur().Pos
		op := p.advance().Literal
		right, err := p.parseCons()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right, Pos: pos}
	}
	return left, nil
}

func (p *Parser) parseCons() (ast.Expr, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	// :: is right-associative
	if p.cur().Type == token.DCOLON {
		pos := p.cur().Pos
		p.advance()
		right, err := p.parseCons()
		if err != nil {
			return nil, err
		}
		return &ast.ConsExpr{Head: left, Tail: right, Pos: pos}, nil
	}
	return left, nil
}

func (p *Parser) parseAddSub() (ast.Expr, error) {
	left, err := p.parseMulDiv()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.PLUS || p.cur().Type == token.MINUS {
		pos := p.cur().Pos
		op := p.advance().Literal
		right, err := p.parseMulDiv()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right, Pos: pos}
	}
	return left, nil
}

func (p *Parser) parseMulDiv() (ast.Expr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.STAR || p.cur().Type == token.SLASH || p.cur().Type == token.PERCENT {
		pos := p.cur().Pos
		op := p.advance().Literal
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &ast.BinOp{Op: op, Left: left, Right: right, Pos: pos}
	}
	return left, nil
}

func (p *Parser) parseUnary() (ast.Expr, error) {
	if p.cur().Type == token.NOT {
		pos := p.cur().Pos
		p.advance()
		expr, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{Op: "!", Expr: expr, Pos: pos}, nil
	}
	if p.cur().Type == token.MINUS {
		// Check if this is a negative number literal vs subtraction
		// Only treat as unary minus if previous token wasn't a value
		pos := p.cur().Pos
		p.advance()
		expr, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryOp{Op: "-", Expr: expr, Pos: pos}, nil
	}
	return p.parseApp()
}

func (p *Parser) parseApp() (ast.Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Function application: f x y == (f(x))(y)
	for p.isAppArg() {
		pos := p.cur().Pos
		arg, err := p.parsePrimary()
		if err != nil {
			return nil, err
		}
		expr = &ast.AppExpr{Func: expr, Arg: arg, Pos: pos}
	}

	return expr, nil
}

func (p *Parser) isAppArg() bool {
	t := p.cur().Type
	return t == token.INT || t == token.STRING || t == token.TRUE || t == token.FALSE ||
		t == token.IDENT || t == token.LPAREN || t == token.LBRACKET
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	switch p.cur().Type {
	case token.INT:
		return p.parseInt()
	case token.TRUE:
		pos := p.cur().Pos
		p.advance()
		return &ast.BoolLit{Value: true, Pos: pos}, nil
	case token.FALSE:
		pos := p.cur().Pos
		p.advance()
		return &ast.BoolLit{Value: false, Pos: pos}, nil
	case token.STRING:
		pos := p.cur().Pos
		val := p.advance().Literal
		return &ast.StringLit{Value: val, Pos: pos}, nil
	case token.IDENT:
		pos := p.cur().Pos
		name := p.advance().Literal
		return &ast.Var{Name: name, Pos: pos}, nil
	case token.LPAREN:
		return p.parseParenOrPair()
	case token.LBRACKET:
		return p.parseList()
	case token.LET:
		return p.parseLetDeclaration()
	case token.IF:
		return p.parseIf()
	case token.FN:
		return p.parseFn()
	case token.FIX:
		return p.parseFix()
	case token.FST:
		return p.parseFst()
	case token.SND:
		return p.parseSnd()
	case token.INL:
		return p.parseInl()
	case token.INR:
		return p.parseInr()
	case token.CASE:
		return p.parseCase()
	case token.PRINT:
		return p.parsePrint(false)
	case token.PRINTLN:
		return p.parsePrint(true)
	default:
		return nil, p.errAt(p.cur().Pos, fmt.Sprintf("unexpected token '%s'", p.cur().Literal))
	}
}

func (p *Parser) parseInt() (ast.Expr, error) {
	pos := p.cur().Pos
	val, err := strconv.Atoi(p.advance().Literal)
	if err != nil {
		return nil, p.errAt(pos, "invalid integer")
	}
	return &ast.IntLit{Value: val, Pos: pos}, nil
}

func (p *Parser) parseParenOrPair() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // (

	// Unit literal
	if p.cur().Type == token.RPAREN {
		p.advance()
		return &ast.UnitLit{Pos: pos}, nil
	}

	first, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	// Pair
	if p.cur().Type == token.COMMA {
		p.advance()
		second, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return &ast.PairExpr{First: first, Second: second, Pos: pos}, nil
	}

	// Parenthesized expression
	if _, err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	return first, nil
}

func (p *Parser) parseList() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // [

	var elems []ast.Expr
	if p.cur().Type != token.RBRACKET {
		elem, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
		for p.cur().Type == token.COMMA {
			p.advance()
			elem, err = p.ParseExpr()
			if err != nil {
				return nil, err
			}
			elems = append(elems, elem)
		}
	}

	if _, err := p.expect(token.RBRACKET); err != nil {
		return nil, err
	}

	// Optional type annotation for empty lists: [] : [Int]
	var elemType ast.Type
	if len(elems) == 0 && p.cur().Type == token.COLON {
		p.advance()
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if lt, ok := t.(*ast.ListType); ok {
			elemType = lt.Elem
		} else {
			elemType = t
		}
	}

	return &ast.ListExpr{Elems: elems, ElemType: elemType, Pos: pos}, nil
}

func (p *Parser) parseIf() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // if

	cond, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(token.THEN); err != nil {
		return nil, err
	}

	then, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(token.ELSE); err != nil {
		return nil, err
	}

	els, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.IfExpr{Cond: cond, Then: then, Else: els, Pos: pos}, nil
}

func (p *Parser) parseFn() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // fn

	if _, err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}

	paramTok, err := p.expect(token.IDENT)
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(token.COLON); err != nil {
		return nil, err
	}

	paramType, err := p.parseType()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}

	if _, err := p.expect(token.FATARROW); err != nil {
		return nil, err
	}

	body, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	return &ast.FnExpr{Param: paramTok.Literal, ParamType: paramType, Body: body, Pos: pos}, nil
}

func (p *Parser) parseFix() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // fix

	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	return &ast.FixExpr{Expr: expr, Pos: pos}, nil
}

func (p *Parser) parseFst() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // fst
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return &ast.FstExpr{Expr: expr, Pos: pos}, nil
}

func (p *Parser) parseSnd() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // snd
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return &ast.SndExpr{Expr: expr, Pos: pos}, nil
}

func (p *Parser) parseInl() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // inl
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	// Require type annotation: inl expr as LeftType + RightType
	var sumType ast.Type
	if p.cur().Type == token.IDENT && p.cur().Literal == "as" {
		p.advance()
		sumType, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	return &ast.InlExpr{Expr: expr, SumType: sumType, Pos: pos}, nil
}

func (p *Parser) parseInr() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // inr
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	var sumType ast.Type
	if p.cur().Type == token.IDENT && p.cur().Literal == "as" {
		p.advance()
		sumType, err = p.parseType()
		if err != nil {
			return nil, err
		}
	}

	return &ast.InrExpr{Expr: expr, SumType: sumType, Pos: pos}, nil
}

func (p *Parser) parseCase() (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // case

	scrutinee, err := p.ParseExpr()
	if err != nil {
		return nil, err
	}

	if _, err := p.expect(token.OF); err != nil {
		return nil, err
	}

	// Optional leading pipe
	if p.cur().Type == token.PIPE {
		p.advance()
	}

	var branches []ast.CaseBranch

	for {
		pat, err := p.parsePattern()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(token.FATARROW); err != nil {
			return nil, err
		}

		body, err := p.ParseExpr()
		if err != nil {
			return nil, err
		}

		branches = append(branches, ast.CaseBranch{Pattern: pat, Body: body})

		if p.cur().Type != token.PIPE {
			break
		}
		p.advance()
	}

	return &ast.CaseExpr{Scrutinee: scrutinee, Branches: branches, Pos: pos}, nil
}

func (p *Parser) parsePattern() (ast.Pattern, error) {
	switch p.cur().Type {
	case token.LBRACKET:
		// [] pattern
		p.advance()
		if _, err := p.expect(token.RBRACKET); err != nil {
			return nil, err
		}
		return &ast.NilPattern{}, nil
	case token.INL:
		p.advance()
		name, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		return &ast.InlPattern{Name: name.Literal}, nil
	case token.INR:
		p.advance()
		name, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		return &ast.InrPattern{Name: name.Literal}, nil
	case token.IDENT:
		// h :: t pattern
		head := p.advance()
		if _, err := p.expect(token.DCOLON); err != nil {
			return nil, err
		}
		tail, err := p.expect(token.IDENT)
		if err != nil {
			return nil, err
		}
		return &ast.ConsPattern{Head: head.Literal, Tail: tail.Literal}, nil
	default:
		return nil, p.errAt(p.cur().Pos, fmt.Sprintf("expected pattern, got '%s'", p.cur().Literal))
	}
}

func (p *Parser) parsePrint(newline bool) (ast.Expr, error) {
	pos := p.cur().Pos
	p.advance() // print / println
	expr, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return &ast.PrintExpr{Expr: expr, Newline: newline, Pos: pos}, nil
}

// Type parsing

func (p *Parser) parseType() (ast.Type, error) {
	return p.parseSumType()
}

func (p *Parser) parseSumType() (ast.Type, error) {
	left, err := p.parseFuncType()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.PLUS {
		p.advance()
		right, err := p.parseFuncType()
		if err != nil {
			return nil, err
		}
		left = &ast.SumType{Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseFuncType() (ast.Type, error) {
	left, err := p.parsePrimaryType()
	if err != nil {
		return nil, err
	}
	// -> is right-associative
	if p.cur().Type == token.ARROW {
		p.advance()
		right, err := p.parseFuncType()
		if err != nil {
			return nil, err
		}
		return &ast.FuncType{Param: left, Return: right}, nil
	}
	return left, nil
}

func (p *Parser) parsePrimaryType() (ast.Type, error) {
	switch p.cur().Type {
	case token.INT_TYPE:
		p.advance()
		return &ast.IntType{}, nil
	case token.BOOL_TYPE:
		p.advance()
		return &ast.BoolType{}, nil
	case token.STRING_TYPE:
		p.advance()
		return &ast.StringType{}, nil
	case token.UNIT_TYPE:
		p.advance()
		return &ast.UnitType{}, nil
	case token.LPAREN:
		p.advance()
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		// Check for pair type
		if p.cur().Type == token.COMMA {
			p.advance()
			t2, err := p.parseType()
			if err != nil {
				return nil, err
			}
			if _, err := p.expect(token.RPAREN); err != nil {
				return nil, err
			}
			return &ast.PairType{First: t, Second: t2}, nil
		}
		if _, err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return t, nil
	case token.LBRACKET:
		p.advance()
		elem, err := p.parseType()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(token.RBRACKET); err != nil {
			return nil, err
		}
		return &ast.ListType{Elem: elem}, nil
	default:
		return nil, p.errAt(p.cur().Pos, fmt.Sprintf("expected type, got '%s'", p.cur().Literal))
	}
}

// ParseSingleExpr parses a single expression (for REPL)
func (p *Parser) ParseSingleExpr() (ast.Expr, error) {
	return p.ParseExpr()
}
