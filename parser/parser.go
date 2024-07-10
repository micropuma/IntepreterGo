package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"fmt"
	"strconv"
)

// ===================== 定义parser的优先级 ================
const (
	_ int = iota

	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X	) 
)

var precedencs = map[token.TokenType] int {
	token.EQ:    EQUALS,
	token.NOT_EQ: EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.SLASH:  PRODUCT,
	token.ASTERISK: PRODUCT,
}

// ====================== 定义parser类 =====================
type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token

	errors []string

	prefixParseFn map[token.TokenType]prefixParseFn
	infixParseFn map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	// 用lexer做初始化parser
	p := &Parser{l : l, 
		         errors : []string{},}		

	// 先读取两个token，已初始化curToken和peekToken
	p.nextToken()
	p.nextToken()

	// 注册expression的优先顺序表
	p.prefixParseFn = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFn = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

// ====================================== helper function ==================================
// 获取下一个token
func (p *Parser)nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()	
}

func (p *Parser)curTokenIs(t token.TokenType) bool{
	return p.curToken.Type == t 
}

func (p *Parser)peekTokenIs(t token.TokenType) bool{
	return p.peekToken.Type == t
}

func (p *Parser)expectPeek(t token.TokenType) bool{
	if p.peekTokenIs(t) {
		// eat up a token and go on
		p.nextToken()
		return true;
	} else {
		p.peekError(t)
		return false;
	}
}

func (p *Parser)peekPrecedence() int {
	if p, ok := precedencs[p.peekToken.Type]; ok {
		return p
	}	

	return LOWEST
}

func (p *Parser)curPrecedence() int {
	if p, ok := precedencs[p.curToken.Type]; ok {
		return p
	}	

	return LOWEST
}

func (p *Parser)Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFn[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFn[tokenType] = fn
}

// =============================================== parse function =======================================
// parse Program ast树
func (p *Parser)ParseProgram() *ast.Program {
	// 初始化空的statement
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// 遍历token，直到EOF
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser)parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFn[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// 在parse infix的时候，需要使用precedence来判断优先级
	for !p.curTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFn[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()	
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser)parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parse Statement ast树结构
func (p *Parser)parseStatement() ast.Statement {
	// 判断是否是let statement
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser)parseLetStatement() *ast.LetStatement {
	// 初始化let statement
	letStmt := &ast.LetStatement{Token: p.curToken}

	// 判断curToken是否是Let
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	
	letStmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	
	// not deal with expression here
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	
	return letStmt
}

func (p *Parser)parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	
	for !p.curTokenIs(token.SEMICOLON) { 
		p.nextToken()
	}

	return returnStmt
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: p.curToken,
		Operator: p.curToken.Literal,
		Left: left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	
	return expression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

