package repl

import (
    "bufio"
    "fmt"
    "io"
    "mitchlang/lexer"
    "mitchlang/token"
)

const Prompt = ">> "

func Start(in io.Reader, out io.Writer) {
    scanner := bufio.NewScanner(in)

    for {
        _, _ = fmt.Fprintf(out, Prompt)
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
