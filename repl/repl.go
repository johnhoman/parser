package repl

import (
	"bufio"
	"fmt"
	"io"
	"mitchlang/object"

	"mitchlang/eval"
	"mitchlang/lexer"
	"mitchlang/parser"
)

var defaultPrompt = ">> "
var prompt = &defaultPrompt

func SetPrompt(s string) {
	*prompt = s
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		_, _ = fmt.Fprintf(out, *prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		obj := eval.Eval(program, env)
		if obj != nil && obj.Type() != object.TypeNull {
			_, _ = io.WriteString(out, obj.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, err := range errors {
		_, _ = io.WriteString(out, "\t"+err+"\n")
	}
}
