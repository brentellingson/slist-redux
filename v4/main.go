package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	expression []rune
	cursor     int
	ast        [][]any
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

func (p *Parser) push(list []any) {
	p.ast = append(p.ast, list)
}

func (p *Parser) pop() []any {
	rval := p.ast[len(p.ast)-1]
	p.ast = p.ast[:len(p.ast)-1]
	return rval
}

func (p *Parser) parseList() []any {
	p.push([]any{}) // parseEntries parses into the top list on the stack
	p.expect('(')
	p.parseEntries()
	p.expect(')')
	return p.pop() // return the top list on the stack
}

func (p *Parser) parseEntries() {
	p.whitespace()

	if p.eof() || p.expression[p.cursor] == ')' {
		return
	}
	entry := p.parseExpression()
	p.ast[len(p.ast)-1] = append(p.ast[len(p.ast)-1], entry)
	p.parseEntries()
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
