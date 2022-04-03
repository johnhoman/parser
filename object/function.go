package object

import (
	"bytes"
	"mitchlang/ast"
	"strings"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

func (f *Function) Type() Type { return TypeFunction }

func (f *Function) Inspect() string {
	out := new(bytes.Buffer)

	params := make([]string, 0)
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Verify implements Object interface
var _ Object = &Function{}
