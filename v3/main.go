package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/brentellingson/slist-redux/v3/expression"
	"github.com/brentellingson/slist-redux/v3/parse"
	"github.com/brentellingson/slist-redux/v3/scan"
	"github.com/brentellingson/slist-redux/v3/token"
)

const TestString = "(defn foo (add 'a 'b))(foo bar)()((foo) bar)\n"

func main() {
	rdr := io.NopCloser(strings.NewReader(TestString))
	defer rdr.Close()

	scanner := scan.NewScanner(rdr)
	tokenizer := token.NewTokenizer(scanner)
	parser := parse.NewParser(tokenizer)
	fmt.Print("Parsing: ", TestString)
	expr := parser.Parse()
	for expr != nil {
		expression.Print(expr)
		fmt.Println()
		expr = parser.Parse()
	}
}
