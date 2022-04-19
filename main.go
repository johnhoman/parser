package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/user"

	"mitchlang/eval"
	"mitchlang/lexer"
	"mitchlang/object"
	"mitchlang/parser"
	"mitchlang/repl"
)

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		raw, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			panic(err)
		}
		script := string(raw)
		l := lexer.New(script)
		p := parser.New(l)
		env := object.NewEnv()
		obj := eval.Eval(p.ParseProgram(), env)
		if obj.Type() == object.TypeError {
			_, _ = io.WriteString(os.Stderr, obj.Inspect())
			_, _ = io.WriteString(os.Stderr, "\n")
			os.Exit(1)
		}
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				_, _ = io.WriteString(os.Stderr, err)
			}
			os.Exit(1)
		}
		os.Exit(0)
	}

	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the MitchLang programming language!\n", u.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
