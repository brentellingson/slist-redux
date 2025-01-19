package expression

import (
	"fmt"
)

type PrintVisitor struct {
	cdr bool
}

func (p *PrintVisitor) VisitAtom(expr *Atom) {
	fmt.Print(expr.Value)
}

func (p *PrintVisitor) VisitPair(expr *Pair) {
	if !p.cdr {
		fmt.Print("(")
	}
	if expr.Car != nil {
		expr.Car.Visit(&PrintVisitor{cdr: false})
	}
	if expr.Car != nil && expr.Cdr != nil {
		fmt.Print(" ")
	}
	if expr.Cdr != nil {
		expr.Cdr.Visit(&PrintVisitor{cdr: true})
	}
	if !p.cdr {
		fmt.Print(")")
	}
}

func Print(expr Expression) {
	p := &PrintVisitor{}
	expr.Visit(p)
}
