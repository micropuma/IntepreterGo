package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// ====================== 定义parser类 =====================
type Parser struct {
	l *lexer.Lexer

	curToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	// 用lexer做初始化parser
	p := &Parser{l : l}		

	// 先读取两个token，已初始化curToken和peekToken
	p.nextToken()
	p.nextToken()

	return p
}

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
		return false;
	}
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

// parse Statement ast树结构
func (p *Parser)parseStatement() ast.Statement {
	// 判断是否是let statement
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil;
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