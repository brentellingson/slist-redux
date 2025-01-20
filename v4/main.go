package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	expression []rune
	cursor     int
}

func Parse(expression string) any {
	p := Parser{expression: []rune(expression)}
	return p.parseExpression()
}

func (p *Parser) eof() bool {
	return p.cursor >= len(p.expression)
}

func (p *Parser) parseExpression() any {
	p.whitespace()
	if p.eof() {
		return nil
	}
	if p.expression[p.cursor] == '(' {
		return p.parseList()
	}
	return p.parseAtom()
}

func (p *Parser) parseList() []any {
	p.expect('(')
	list := p.parseEntries([]any{})
	p.expect(')')
	return list
}

func (p *Parser) parseEntries(list []any) []any {
	p.whitespace()
	if p.eof() || p.expression[p.cursor] == ')' {
		return list
	}
	entry := p.parseExpression()
	list = append(list, entry)
	return p.parseEntries(list)
}

func (p *Parser) parseAtom() any {
	var atom []rune
	for !p.eof() && !iswhitespace(p.expression[p.cursor]) && p.expression[p.cursor] != ')' {
		atom = append(atom, p.expression[p.cursor])
		p.cursor++
	}
	if ival, err := strconv.Atoi(string(atom)); err == nil {
		return ival
	}
	return string(atom)
}

func (p *Parser) expect(c rune) {
	if p.eof() {
		panic("unexpected end of input, expected " + string(c))
	}
	if p.expression[p.cursor] != c {
		panic("unexpected character " + string(p.expression[p.cursor]) + ", expected " + string(c))
	}
	p.cursor++
}

func iswhitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (p *Parser) whitespace() {
	for !p.eof() && iswhitespace(p.expression[p.cursor]) {
		p.cursor++
	}
}

func print(expr string) {
	fmt.Println(expr, "->", Parse(expr))
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
