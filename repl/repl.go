package repl

import (
	"bufio"
	"fmt"
	"io"

	"mitchlang/lexer"
	"mitchlang/token"
)

var defaultPrompt = ">> "
var prompt *string = &defaultPrompt

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
		lxr := lexer.New(line)
		for tok := lxr.NextToken(); tok.Type != token.EOF; tok = lxr.NextToken() {
			_, _ = fmt.Fprintf(out, "%+v\n", tok)
		}
	}

}
