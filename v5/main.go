package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	reader  io.RuneReader
	current rune
	eof     bool
}

func Parse(reader io.Reader) any {
	p := Parser{reader: bufio.NewReader(reader)}
	p.advance() // prime the pump!
	return p.parseExpression()
}

func (p *Parser) advance() {
	r, _, err := p.reader.ReadRune()
	if err == io.EOF {
		p.eof = true
	} else if err != nil {
		panic("error reading stream " + err.Error())
	}
	p.current = r
}

func (p *Parser) parseExpression() any {
	p.whitespace()
	switch {
	case p.eof:
		return nil
	case p.current == '(':
		return p.parseList()
	default:
		return p.parseAtom()
	}
}

func (p *Parser) parseList() []any {
	list := []any{}
	p.expect('(')
	p.whitespace()
	for !p.eof && p.current != ')' {
		list = append(list, p.parseExpression())
		p.whitespace()
	}
	p.expect(')')
	return list
}

func (p *Parser) parseAtom() any {
	var sval []rune
	for !p.eof && !iswhitespace(p.current) && p.current != ')' {
		sval = append(sval, p.current)
		p.advance()
	}
	if ival, err := strconv.Atoi(string(sval)); err == nil {
		return ival
	}
	return string(sval)
}

func (p *Parser) expect(c rune) {
	if p.eof {
		panic("unexpected end of input, expected " + string(c))
	}
	if p.current != c {
		panic("unexpected character " + string(p.current) + ", expected " + string(c))
	}
	p.advance()
}

func iswhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (p *Parser) whitespace() {
	for !p.eof && iswhitespace(p.current) {
		p.advance()
	}
}

func print(expr string) {
	fmt.Println(expr, "->", Parse(strings.NewReader(expr)))
}

func main() {
	// nil
	print(``)

	// atoms
	print(`a`)
	print(`foo`)

	// empty lists
	print(`()`)
	print(`( )`)
	print(`( ( ) )`)
	print(`(()(()))`)

	// non-empty lists
	print(`(+ 1 2)`)
	print(`(* (+ 1 2) (- 3 5))`)

	// multi-line
	print(`
		(define (fact n)
		  (if (= n 0)
		    1
			(* n (fact (- n 1)))))`)
}
