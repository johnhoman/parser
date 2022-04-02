package repl

import (
	"bufio"
	"fmt"
	"io"

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
		obj := eval.Eval(program)
		if obj != nil {
			_, _ = io.WriteString(out, obj.Inspect())
			_, _ = io.WriteString(out, "\n")
		}
	}

}

func printParserErrors(out io.Writer, errors []string) {
	_, _ = fmt.Fprintln(
		out,
		"Whoops! You fucked something up and I don't seem to know how to fix it",
	)
	for _, err := range errors {
		_, _ = io.WriteString(out, "\t"+err+"\n")
	}
}
