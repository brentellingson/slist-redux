package parse

import (
	"github.com/brentellingson/slist-redux/v3/expression"
	"github.com/brentellingson/slist-redux/v3/token"
)

type Parser struct {
	tokenizer *token.Tokenizer
}

func NewParser(tokenizer *token.Tokenizer) *Parser {
	tokenizer.NextToken() // prime the pump
	return &Parser{tokenizer}
}

func (p *Parser) Parse() expression.Expression {
	if p.tokenizer.Current.Kind == token.EOF {
		return nil
	}

	if p.tokenizer.Current.Kind == token.OpenParen {
		return p.ParseList()
	}

	return p.ParseAtom()
}

func (p *Parser) ParseList() *expression.Pair {
	head := &expression.Pair{}
	tail := head
	p.tokenizer.NextToken() // consume '('
	for {
		if p.tokenizer.Current.Kind == token.EOF {
			panic("unexpected EOF")
		}
		if p.tokenizer.Current.Kind == token.CloseParen {
			p.tokenizer.NextToken() // consume ')' in ()
			break
		}
		tail.Car = p.Parse()
		if p.tokenizer.Current.Kind == token.CloseParen {
			p.tokenizer.NextToken() // consume ')' in (a)
			break
		}
		tail.Cdr = &expression.Pair{}
		tail = tail.Cdr.(*expression.Pair)
	}
	return head
}

func (p *Parser) ParseAtom() *expression.Atom {
	retval := &expression.Atom{Value: p.tokenizer.Current.Value}
	p.tokenizer.NextToken()
	return retval
}
