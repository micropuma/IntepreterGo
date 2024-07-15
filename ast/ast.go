package ast

import ("monkey/token"
		"bytes"
	   )

// ============================== 接口定义 ============================== //
// Node接口必须实现TokenLiteral
type Node interface {
	TokenLiteral() string
	String() string
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

type ReturnStatement struct {
	Token token.Token   // return token
	ReturnValue Expression
}

type ExpressionStatement struct {
	Token token.Token   
	Expression Expression
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type PrefixExpression struct {
	Token token.Token      // 记录比如！这种
	Operator string
	Right Expression
}

type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

type Boolean struct {
	Token token.Token
	Value bool
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

// fulfill interface
func (rs *ReturnStatement)statementNode() {}
func (rs *ReturnStatement)TokenLiteral() string { return rs.Token.Literal }

// fulfill interface
func (es *ExpressionStatement)statementNode() {}
func (es *ExpressionStatement)TokenLiteral() string { return es.Token.Literal } 

func (p *Program)String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (ls *LetStatement)String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

func (rs *ReturnStatement)String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

func (es *ExpressionStatement)String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (i *Identifier)String() string {
	return i.Value
}

func (il *IntegerLiteral) expressionNode()   {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string { return il.Token.Literal }

func (pl *PrefixExpression) expressionNode()  {}
func (pl *PrefixExpression) TokenLiteral() string { return pl.Token.Literal }
func (pl *PrefixExpression) String() string {
	var out bytes.Buffer
	
	out.WriteString("(")
	out.WriteString(pl.Operator)
	out.WriteString(pl.Right.String())
	out.WriteString(")")

	return out.String()
}

func (il *InfixExpression) expressionNode()  {}
func (il *InfixExpression) TokenLiteral() string { return il.Token.Literal }
func (il *InfixExpression) String() string {
	var out bytes.Buffer
	
	out.WriteString("(")
	out.WriteString(il.Left.String())
	out.WriteString(" " + il.Operator + " ")
	out.WriteString(il.Right.String())
	out.WriteString(")")

	return out.String()
}

func (b *Boolean) expressionNode()             {}
func (b *Boolean) TokenLiteral() string        { return b.Token.Literal }
func (b *Boolean) String() string              { return b.Token.Literal }