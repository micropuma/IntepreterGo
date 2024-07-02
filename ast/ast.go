package ast

import "monkey/token"

// ============================== 接口定义 ============================== //
// Node接口必须实现TokenLiteral
type Node interface {
	TokenLiteral() string
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

// ==================== AST 结构定义 ======================== //
type Program struct {
	Statements []Statement
}

func (p *Program)TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	Token token.Token   // Let token
	Name *Identifier
	Value Expression
}

// fulfill interface
func (l *LetStatement)statementNode() {}
func (l *LetStatement)TokenLiteral() string { return l.Token.Literal }

type Identifier struct {
	Token token.Token     // IDENT token
	Value string
}

// fulfill interface
func (i *Identifier)expressionNode() {}
func (i *Identifier)TokenLiteral() string { return i.Token.Literal }
