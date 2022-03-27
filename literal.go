package parser

type LiteralType string
const (
    StringLiteralType LiteralType = "STRING_LITERAL"
    IntLiteralType    LiteralType = "INT_LITERAL"
    WhitespaceLiteral LiteralType = "WHITESPACE_LITERAL"
    CommentLiteral    LiteralType = "COMMENT_LITERAL"
    ExpressionTerm    LiteralType = "EXPRESSION_TERM"
)


type Literal interface {
    Type() LiteralType
    Value() interface{}
}

type IntLiteral struct {
    value int
}

func (lit *IntLiteral) Value() interface{} {
    return lit.value
}

func (lit *IntLiteral) Type() LiteralType {
    return IntLiteralType
}

type StringLiteral struct {
    value string
}

func (lit *StringLiteral) Value() interface{} {
    return lit.value
}

func (lit *StringLiteral) Type() LiteralType {
    return StringLiteralType
}

func repr(literal Literal) map[string]interface{} {
    return map[string]interface{}{
        "Type": literal.Type(),
        "Value": literal.Value(),
    }

}