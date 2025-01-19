package token

import (
	scan "github.com/brentellingson/slist-redux/v3/scan"
)

//go:generate stringer -type=TokenKind
type TokenKind uint8

const (
	EOF TokenKind = iota
	OpenParen
	CloseParen
	Symbol
)

type Token struct {
	Kind  TokenKind
	Value string
}

func whitespace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

type Tokenizer struct {
	scanner *scan.Scanner
	buffer  []rune
	Current Token
}

func NewTokenizer(scanner *scan.Scanner) *Tokenizer {
	scanner.Advance() // prime the pump
	return &Tokenizer{scanner: scanner}
}

func (t *Tokenizer) NextToken() {
	t.Current = t.nextToken()
}

// NextToken returns the next token from the scanner.
// it expects to to be called with scanner.Current pointing at the first rune in the token,
// after its called, scanner.Current will point at the first run after the token (and the presumably the next token)
func (t *Tokenizer) nextToken() Token {
	// skip whitespace!
	for !t.scanner.EOF && whitespace(t.scanner.Current) {
		t.scanner.Advance()
	}

	if t.scanner.EOF {
		return Token{Kind: EOF, Value: ""}
	}

	switch t.scanner.Current {
	case '(':
		t.scanner.Advance()
		return Token{Kind: OpenParen, Value: "("}
	case ')':
		t.scanner.Advance()
		return Token{Kind: CloseParen, Value: ")"}
	default:
		return t.nextAtom()
	}
}

func (t *Tokenizer) nextAtom() Token {
	t.buffer = t.buffer[:0]

	for {
		t.buffer = append(t.buffer, t.scanner.Current)
		t.scanner.Advance()
		if t.scanner.EOF || whitespace(t.scanner.Current) || t.scanner.Current == ')' {
			break
		}
	}
	return Token{Kind: Symbol, Value: string(t.buffer)}
}
