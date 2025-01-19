package expression

type Visitor interface {
	VisitAtom(expr *Atom)
	VisitPair(expr *Pair)
}

type Expression interface {
	Visit(visitor Visitor)
}

type Atom struct {
	Value string
}

func (expr *Atom) Visit(visitor Visitor) {
	visitor.VisitAtom(expr)
}

type Pair struct {
	Car Expression
	Cdr Expression
}

func (expr *Pair) Visit(visitor Visitor) {
	visitor.VisitPair(expr)
}
