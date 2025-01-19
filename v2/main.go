package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type char uint8

const TestString = "(defn foo (add 'a 'b)) (foo bar)\n"

type Atom struct {
	issymbol bool
	symbol   string
	number   int
}

type Pair struct {
	car *SExp
	cdr *SExp
}

type SExp struct {
	isatom bool
	atom   Atom
	pair   Pair
}

func (sexp *SExp) Print() {
	sexp.PrintOne(true)
}

func (sexp *SExp) PrintOne(doparen bool) {
	if sexp == nil {
		return
	}
	if sexp.isatom {
		if sexp.atom.issymbol {
			fmt.Print(sexp.atom.symbol)
		} else {
			fmt.Print(sexp.atom.number)
		}
		return
	}

	if doparen {
		fmt.Print("(")
	}
	sexp.pair.car.PrintOne(true)
	if sexp.pair.car != nil && sexp.pair.cdr != nil {
		fmt.Print(" ")
	}
	sexp.pair.cdr.PrintOne(false)
	if doparen {
		fmt.Print(")")
	}
}

const EOF = -1

type Token rune // ')', '(', 'A', '0', EOF
var (
	token    Token                  // current token
	tokenbuf = make([]rune, 0, 100) // contents of current token
)

func main() {
	scanner := OpenFile()
	NextToken(scanner) // prime the pump
	for {
		sexp := Parse(scanner)
		if sexp == nil {
			break
		}
		sexp.Print()
		fmt.Println()
	}
}

func OpenFile() io.RuneScanner {
	return bufio.NewReader(strings.NewReader(TestString))
}

func WhiteSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func NextToken(scanner io.RuneScanner) {
	// skip whitspace
	var c rune
	var err error
	for {
		c, _, err = scanner.ReadRune()
		if err == io.EOF {
			token = EOF
			return
		}
		if err != nil {
			panic(err)
		}
		if !WhiteSpace(c) {
			break
		}
	}

	switch c {
	case '(', ')':
		token = Token(c)
	default:
		tokenbuf = tokenbuf[:0]

		for {
			tokenbuf = append(tokenbuf, c)
			c, _, err = scanner.ReadRune()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			if WhiteSpace(c) || c == ')' {
				err := scanner.UnreadRune()
				if err != nil {
					panic(err)
				}
				break
			}
		}
		if tokenbuf[0] >= '0' && tokenbuf[0] <= '9' {
			token = '0'
		} else {
			token = 'A'
		}
	}
}

func ParseList(scanner io.RuneScanner) *SExp {
	NextToken(scanner)
	sexp := &SExp{pair: Pair{}}
	retval := sexp
	for {
		if token == EOF {
			panic("unexpected EOF")
		}
		if token == ')' {
			break
		}
		sexp.pair.car = Parse(scanner)
		if token == ')' {
			break
		}

		sexp.pair.cdr = &SExp{pair: Pair{}}
		sexp = sexp.pair.cdr
	}
	return retval
}

func Parse(scanner io.RuneScanner) *SExp {
	var sexp *SExp

	switch token {
	case EOF:
		return nil
	case '(':
		sexp = ParseList(scanner)
	case '0':
		val, _ := strconv.Atoi(string(tokenbuf))
		sexp = &SExp{isatom: true, atom: Atom{issymbol: false, number: val}}
	case 'A':
		sexp = &SExp{isatom: true, atom: Atom{issymbol: true, symbol: string(tokenbuf)}}
	default:
		panic("unknown token " + string(token))
	}
	NextToken(scanner)
	return sexp
}
