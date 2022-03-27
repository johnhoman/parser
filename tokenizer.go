package parser

import (
	"fmt"
	"regexp"
	"strconv"
)

type TokenType int
type BinaryOp  int

const (
	IntegerLiteralRegex      = `(^\d+)`
	StringLiteralRegex       = `^"([^"]*)"`
	WhitespaceRegex          = `^\s+`
	CommentRegex             = `^//.*`
	BlockCommentRegex        = `^/\*[\s\S]*?\*/`
	TerminateExpressionRegex = `^;`
	BlockStatementBeginRegex = `^{`
	BlockStatementEndRegex   = `^}`
    BinaryOpRegex            = `^(\+|\-)`
    BinaryOpAddRegex         = `^(\+)`

	IntegerLiteralToken TokenType = iota + 1
	StringLiteralToken
	NullToken
	TerminateLiteralToken
	BlockStatementBeginToken
	BlockStatementEndToken
    BinaryOpAddToken

    BinaryOpAdd BinaryOp = iota + 1
)

var binaryOperations = map[TokenType]BinaryOp{
    BinaryOpAddToken: BinaryOpAdd,
}

func isBinaryOp(tokenType TokenType) bool {
    if _, ok := binaryOperations[tokenType]; ok {
        return true
    }
    return false
}

func (tt TokenType) String() string {
	switch tt {
	case IntegerLiteralToken:
		return "IntegerLiteral"
	case StringLiteralToken:
		return "StringLiteral"
	case NullToken:
		return "Null"
	case TerminateLiteralToken:
		return "TerminateExpression"
	case BlockStatementBeginToken:
		return "BlockStatementBegin"
	case BlockStatementEndToken:
		return "BlockStatementEnd"
	default:
		return fmt.Sprintf("%d", int(tt))
	}
}

var specs = map[string]TokenType{
	IntegerLiteralRegex:      IntegerLiteralToken,
	StringLiteralRegex:       StringLiteralToken,
	WhitespaceRegex:          NullToken,
	CommentRegex:             NullToken,
	BlockCommentRegex:        NullToken,
	TerminateExpressionRegex: TerminateLiteralToken,
	BlockStatementBeginRegex: BlockStatementBeginToken,
	BlockStatementEndRegex:   BlockStatementEndToken,
    BinaryOpAddRegex:         BinaryOpAddToken,
}

// Lexical analysis

type SyntaxError struct {
	message string
}

func (err *SyntaxError) Error() string {
	return err.message
}

func NewSyntaxError(message string) error {
	return &SyntaxError{message}
}

type BinaryExpression struct {
	Operator BinaryOp
	Left     *Expression
	Right    *Expression
}

type Expression struct {
	Literal          *Literal
	BinaryExpression *BinaryExpression
}

type ExpressionStatement struct {
	Expression Expression
}

type BlockStatement struct {
	StatementList StatementList
}

type EmptyStatement struct{}

type Statement struct {
	ExpressionStatement *ExpressionStatement
	BlockStatement      *BlockStatement
	EmptyStatement      *EmptyStatement
}

type StatementList struct {
	Statements    []Statement
	StatementList *StatementList
}

type Program struct {
	Type string
	Body StatementList
}

type Parser interface {
	// Parse parses a string into an abstract syntax tree (AST)
	Parse(s string) (Program, error)

	// Program
	//   ; StatementList
	//   ;
	Program() (Program, error)

	// StatementList
	//   : Statement
	//   | StatementList Statement
	//   ;
	StatementList(TokenType) (StatementList, error)

	// Statement
	//   : ExpressionStatement
	//   | BlockStatement
	//   ;
	Statement() (Statement, error)

	// ExpressionStatement
	//   : Expression ;
	//   ;
	ExpressionStatement() (Statement, error)

	// BlockStatement
	//   : '{' StatementList '}'
	//   ;
	BlockStatement() (Statement, error)

    // AdditiveExpression
    //   : Literal
    //   | AdditiveExpression ADDITIVE_EXPRESSION Literal
    //   ;

    AdditiveExpression() (Expression, error)

	// Expression
	//   : Literal
	//   ;
	Expression() (Expression, error)

	// Literal
	//   : StringLiteral
	//   | NumericLiteral
	//   ;
	Literal() (Literal, error)
}

type parser struct {
	tokenizer *tokenizer
	lookAhead Token
}

// Parse the string s into an abstract syntax tree
func (p *parser) Parse(s string) (Program, error) {
	p.tokenizer = NewTokenizer(s)
	p.lookAhead, _ = p.tokenizer.NextToken()
	return p.Program()
}

func (p *parser) Expression() (Expression, error) {
    expr, err := p.AdditiveExpression()
    if err != nil {
        return Expression{}, err
    }
	_, err = p.eat(TerminateLiteralToken)
	if err != nil {
		return Expression{}, err
	}
	return expr, nil
}

func (p *parser) ExpressionStatement() (Statement, error) {
	expression, err := p.Expression()
	if err != nil {
		return Statement{}, err
	}
	return Statement{
		ExpressionStatement: &ExpressionStatement{Expression: expression},
	}, nil
}

func (p *parser) BlockStatement() (Statement, error) {
	_, err := p.eat(BlockStatementBeginToken)
	if err != nil {
		return Statement{}, err
	}
	block := &BlockStatement{}
	if p.lookAhead.Type != BlockStatementEndToken {
		var err error
		block.StatementList, err = p.StatementList(BlockStatementEndToken)
		if err != nil {
			return Statement{}, err
		}
	} else {
		block.StatementList = StatementList{}
	}
	_, err = p.eat(BlockStatementEndToken)
	if err != nil {
		return Statement{}, err
	}
	return Statement{BlockStatement: block}, nil
}

func (p *parser) EmptyStatement() (Statement, error) {
	_, err := p.eat(TerminateLiteralToken)
	if err != nil {
		return Statement{}, err
	}
	return Statement{EmptyStatement: &EmptyStatement{}}, nil
}

func (p *parser) Statement() (Statement, error) {
	switch p.lookAhead.Type {
	case BlockStatementBeginToken:
		return p.BlockStatement()
	case TerminateLiteralToken:
		return p.EmptyStatement()
	default:
		return p.ExpressionStatement()
	}
}

func (p *parser) StatementList(stop TokenType) (StatementList, error) {
	statements := StatementList{}

	for !p.lookAhead.IsEmpty() && p.lookAhead.Type != stop {
		statement, err := p.Statement()
		if err != nil {
			return StatementList{}, err
		}
		statements.Statements = append(statements.Statements, statement)
	}
	return statements, nil
}

func (p *parser) Program() (Program, error) {
	statements, err := p.StatementList(TokenType(0))
	if err != nil {
		return Program{}, err
	}
	return Program{Type: "Program", Body: statements}, nil
}

func (p *parser) eat(tokenType TokenType) (Token, error) {
	token := p.lookAhead
	if token.IsEmpty() {
		return token, NewSyntaxError("EOF")
	}
	if token.Type != tokenType {
		return token, NewSyntaxError(fmt.Sprintf("unexpected token: '%s'", token.Value))
	}
	var err error
	p.lookAhead, err = p.tokenizer.NextToken()
	if err != nil {
		return Token{}, err
	}
	return token, nil
}

func (p *parser) IntLiteral() (Literal, error) {
	token, err := p.eat(IntegerLiteralToken)
	if err != nil {
		return Literal{}, err
	}
	i, _ := strconv.Atoi(token.Value)
	return Literal{Value: i}, nil
}

func (p *parser) StringLiteral() (Literal, error) {
	token, err := p.eat(StringLiteralToken)
	if err != nil {
		return Literal{}, err
	}
	return Literal{Value: token.Value}, nil
}

func (p *parser) Literal() (Literal, error) {
	switch p.lookAhead.Type {
	case StringLiteralToken:
		{
			return p.StringLiteral()
		}
	case IntegerLiteralToken:
		{
			return p.IntLiteral()
		}
	default:
		{
			return Literal{}, NewSyntaxError(fmt.Sprintf("Invalid literal type %s", p.lookAhead.Type))
		}
	}
}

func (p *parser) AdditiveExpression() (Expression, error) {
    lit, err := p.Literal()
    if err != nil {
        return Expression{}, err
    }
    left := Expression{Literal: &lit}

    for isBinaryOp(p.lookAhead.Type) {
        op, err := p.eat(p.lookAhead.Type)
        if err != nil {
            return Expression{}, err
        }
        right, err := p.Literal()
        if err != nil {
            return Expression{}, err
        }
        left = Expression{BinaryExpression: &BinaryExpression{
            Operator: binaryOperations[op.Type],
            Left: &Expression{
                BinaryExpression: left.BinaryExpression,
                Literal: left.Literal,
            },
            Right: &Expression{Literal: &right},
        }}
    }
    return left, nil
}

var _ Parser = &parser{}

func New() *parser {
	return &parser{}
}

type Token struct {
	Type  TokenType
	Value string
}

func (tok *Token) IsEmpty() bool {
	return tok.Type == 0
}

type String string

func (s String) Len() int {
	return len(s)
}

func (s String) Slice(start int) String {
	return s[start:]
}

type tokenizer struct {
	String
	cursor int
}

func (tok *tokenizer) hasMoreTokens() bool { return tok.cursor < tok.String.Len() }

func (tok *tokenizer) NextToken() (Token, error) {
	skip := map[TokenType]bool{
		TerminateLiteralToken:    false,
		BlockStatementBeginToken: false,
		BlockStatementEndToken:   false,
	}
	if !tok.hasMoreTokens() {
		return Token{}, nil
	}
	for pattern, literalType := range specs {
		str := string(tok.String.Slice(tok.cursor))

		re := regexp.MustCompile(pattern)
		if re.MatchString(str) {
			match := re.FindStringSubmatch(str)
			tok.cursor += len(match[0])
			if literalType == NullToken {
				return tok.NextToken()
			}
			if _, ok := skip[literalType]; ok {
				return Token{Type: literalType}, nil
			}
			return Token{Type: literalType, Value: match[1]}, nil
		}
	}
	return Token{}, NewSyntaxError(fmt.Sprintf(`Unexpected token: %c`, tok.String[0]))
}

func NewTokenizer(s string) *tokenizer {
	return &tokenizer{String(s), 0}
}
