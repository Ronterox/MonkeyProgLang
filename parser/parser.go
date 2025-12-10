package parser

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
	"strings"
)

const (
	_ = iota
	LOWEST
	OR
	AND
	EQUALS
	LESSGREATER // < or >
	SUM
	PRODUCT
	PREFIX
	CALL
	INDEX
)

var precedences = map[token.TokenType]int{
	token.OR:       OR,
	token.AND:      AND,
	token.EQ:       EQUALS,
	token.NE:       EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LE:       LESSGREATER,
	token.GE:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.PERCENT:  PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{l: lexer}
	// Read 2 tokens so that currToken and peekToken are set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.IDENT] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.TRUE] = p.parseBooleanExpression
	p.prefixParseFns[token.FALSE] = p.parseBooleanExpression
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression
	p.prefixParseFns[token.IF] = p.parseIfElseExpression
	p.prefixParseFns[token.FUNCTION] = p.parseFunctionExpression
	p.prefixParseFns[token.STRING] = p.parseString
	p.prefixParseFns[token.TEMPLATE] = p.parseTemplate
	p.prefixParseFns[token.LBRACKET] = p.parseArrayLiteral
	p.prefixParseFns[token.LBRACE] = p.parseHashLiteral

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.PERCENT] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.AND] = p.parseInfixExpression
	p.infixParseFns[token.OR] = p.parseInfixExpression
	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NE] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.LE] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.GE] = p.parseInfixExpression
	p.infixParseFns[token.LPAREN] = p.parseCallExpression
	p.infixParseFns[token.LBRACKET] = p.parseIndexExpression
	return p
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		smtm := p.parseStatement()
		if smtm != nil {
			program.Statements = append(program.Statements, smtm)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfElseExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.currToken}

	if p.peekTokenIs(token.LPAREN) {
		p.nextToken()
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.ParseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.ParseBlockStatement()
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
	exp := &ast.HashLiteral{Token: p.currToken, Pairs: map[ast.Expression]ast.Expression{}}

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.nextToken()

		value := p.parseExpression(LOWEST)
		exp.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	exp := &ast.ArrayLiteral{ExpressionsContainer: ast.ExpressionsContainer{Token: p.currToken}}
	p.nextToken()

	if !p.currTokenIs(token.RBRACKET) {
		param := p.parseExpression(LOWEST)
		exp.Elements = append(exp.Elements, param)

		for p.peekTokenIs(token.COMMA) {
			p.nextToken()

			if p.peekTokenIs(token.RBRACKET) {
				break
			}

			p.nextToken()
			param := p.parseExpression(LOWEST)
			exp.Elements = append(exp.Elements, param)
		}

		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
	}

	return exp
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	exp := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()

	if !p.currTokenIs(token.RPAREN) {
		param := p.parseIdentifier().(*ast.Identifier)
		exp.Parameters = append(exp.Parameters, param)

		for p.peekTokenIs(token.COMMA) {
			p.nextToken()
			p.nextToken()
			param := p.parseIdentifier().(*ast.Identifier)
			exp.Parameters = append(exp.Parameters, param)
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Body = p.ParseBlockStatement()
	return exp
}

func (p *Parser) ParseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseString() ast.Expression {
	var literal bytes.Buffer
	literal.WriteString(p.currToken.Literal)
	for p.peekTokenIs(token.STRING) {
		p.nextToken()
		literal.WriteString(p.currToken.Literal)
	}
	preprocessed := strings.ReplaceAll(literal.String(), "\\n", "\n")
	return &ast.StringLiteral{Token: p.currToken, Value: preprocessed}
}

func (p *Parser) parseTemplate() ast.Expression {
	tmpl := &ast.TemplateString{ExpressionsContainer: ast.ExpressionsContainer{Token: p.currToken}}
	tmpl.Elements = append(tmpl.Elements, &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal})
	for p.peekTokenIs(token.TEMPLATE) || p.peekTokenIs(token.IDENT) {
		var exp ast.Expression

		p.nextToken()
		if p.currTokenIs(token.TEMPLATE) {
			exp = &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
		} else {
			exp = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		}

		tmpl.Elements = append(tmpl.Elements, exp)
	}
	return tmpl
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// defer untrace(trace("parsePrefixExpression"))
	pe := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}
	p.nextToken()
	pe.Right = p.parseExpression(PREFIX)
	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
	pe := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}
	precendence := p.currPrecedence()
	p.nextToken()
	pe.Right = p.parseExpression(precendence)

	return pe
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseInfixExpression"))
	pe := &ast.IndexExpression{
		Token: p.currToken,
		Left:  left,
	}

	p.nextToken()

	pe.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return pe
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	// defer untrace(trace("parseCallExpression"))
	ce := &ast.CallExpression{Token: p.currToken, Function: left}

	p.nextToken()
	if p.currTokenIs(token.RPAREN) {
		return ce
	}

	ce.Arguments = append(ce.Arguments, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ce.Arguments = append(ce.Arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return ce
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// defer untrace(trace("parseExpressionStatement"))
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// defer untrace(trace("parseExpression"))
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		msg := fmt.Sprintf("no prefix parse function for %s found", p.currToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	leftExp := prefix()

	for precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmt.RetValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// Only advances to the next token if is the expected one
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) currTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}
