package ast

import (
	"bytes"
	"mitchlang/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	out := new(bytes.Buffer)
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type LetStatement struct {
	Token *token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	out := new(bytes.Buffer)

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type Identifier struct {
	Token *token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token       *token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	out := new(bytes.Buffer)

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      *token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token *token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    *token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	out := new(bytes.Buffer)

	out.WriteByte('(')
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteByte(')')

	return out.String()
}

type InfixExpression struct {
	Token    *token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	out := new(bytes.Buffer)

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

type Boolean struct {
	Token *token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       *token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	out := new(bytes.Buffer)

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type BlockStatement struct {
	Token      *token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	out := new(bytes.Buffer)

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionLiteralExpression struct {
	Token      *token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fle *FunctionLiteralExpression) expressionNode()      {}
func (fle *FunctionLiteralExpression) TokenLiteral() string { return fle.Token.Literal }
func (fle *FunctionLiteralExpression) String() string {
	out := new(bytes.Buffer)

	params := make([]string, 0, len(fle.Parameters))
	for _, param := range fle.Parameters {
		params = append(params, param.String())
	}

	out.WriteString(fle.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(" { ")
	out.WriteString(fle.Body.String())
	out.WriteString(" }")
	return out.String()
}

type CallExpression struct {
	Token     *token.Token
	Function  Expression // Identifier || FunctionLiteralExpression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	out := new(bytes.Buffer)

	args := make([]string, 0)
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type StringLiteral struct {
	Token *token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string {
	out := new(bytes.Buffer)
	out.WriteByte('"')
	out.WriteString(sl.Token.Literal)
	out.WriteByte('"')
	return out.String()
}

type ListExpression struct {
	Token *token.Token
	Items []Expression
}

func (sl *ListExpression) expressionNode()      {}
func (sl *ListExpression) TokenLiteral() string { return sl.Token.Literal }
func (sl *ListExpression) String() string {

	items := make([]string, 0, len(sl.Items))
	for _, item := range sl.Items {
		items = append(items, item.String())
	}
	out := new(bytes.Buffer)
	out.WriteByte('[')
	out.WriteString(strings.Join(items, ", "))
	out.WriteByte(']')
	return out.String()
}

type IndexExpression struct {
	Token *token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	out := new(bytes.Buffer)

	out.WriteByte('(')
	out.WriteString(ie.Left.String())
	out.WriteByte('[')
	out.WriteString(ie.Index.String())
	out.WriteByte(']')
	out.WriteByte(')')
	return out.String()
}

type mapEntry struct {
	key   string
	value string
}

type MapExpression struct {
	Token   *token.Token
	Entries map[Expression]Expression
}

func (exp *MapExpression) expressionNode()      {}
func (exp *MapExpression) TokenLiteral() string { return exp.Token.Literal }
func (exp *MapExpression) String() string {
	out := new(bytes.Buffer)

	pairs := make([]mapEntry, 0, len(exp.Entries))
	for key, value := range exp.Entries {
		pairs = append(pairs, mapEntry{key.String(), value.String()})
	}

	out.WriteByte('{')
	for k, pair := range pairs {
		out.WriteString(pair.key + ": " + pair.value)
		if k < len(pairs) {
			out.WriteString(", ")
		}
	}
	out.WriteByte('}')
	return out.String()
}
