package parser_test

import (
    "encoding/json"
    "fmt"
    "testing"

    "parser"
)

func TestParser(t *testing.T) {
    ast := parser.New().Parse(`42`)

    encoded, _ := json.MarshalIndent(ast, "", "  ")
    fmt.Printf("%s\n", encoded)
}