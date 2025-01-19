package main

import "fmt"

type char uint8

const TestString = "(defn foo (add 'a 'b))\n"

type Atom struct {
	string  *[100]char
	integer int
	next    *Slist
}

type List struct {
	car *Slist
	cdr *Slist
}

type Slist struct {
	isatom   bool
	isstring bool
	atom     Atom
	list     List
}

func (this *Slist) Car() *Slist {
	return this.list.car
}

func (this *Slist) Cdr() *Slist {
	return this.list.cdr
}

func (this *Slist) String() *[100]char {
	return this.atom.string
}

func (this *Slist) Integer() int {
	return this.atom.integer
}

func (slist *Slist) Free() {
	if slist == nil {
		return
	}
	if slist.isatom {
	} else {
		slist.Car().Free()
		slist.Cdr().Free()
	}
}

func (slist *Slist) Print() {
	slist.PrintOne(true)
	fmt.Print("\n")
}

func (slist *Slist) PrintOne(doparen bool) {
	if slist == nil {
		return
	}
	if slist.isatom {
		if slist.isstring {
			for i := 0; (*slist.atom.string)[i] != 0; i++ {
				fmt.Print(string((*slist.atom.string)[i]))
			}
		} else {
			fmt.Print(slist.Integer())
		}
	} else {
		if doparen {
			fmt.Print("(")
		}
		slist.Car().PrintOne(true)
		if slist.Cdr() != nil {
			print(" ")
			slist.Cdr().PrintOne(false)
		}
		if doparen {
			fmt.Print(")")
		}
	}
}

var (
	token      int
	peekc      int
	lineno     int = 1
	input      [100 * 1000]char
	inputindex int
	tokenbuf   [100]char
)

const EOF = -1

func main() {
	var list *Slist
	OpenFile()
	for {
		list = Parse()
		if list == nil {
			break
		}
		list.Print()
		list.Free()
		break
	}
}

func Get() int {
	var c int
	if peekc >= 0 {
		c = peekc
		peekc = -1
	} else {
		c = int(input[inputindex])
		inputindex++
		if c == '\n' {
			lineno++
		}
		if c == 0 {
			inputindex--
			c = EOF
		}
	}
	return c
}

func WhiteSpace(c int) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func NextToken() {
	var i, c int
	tokenbuf[0] = 0
	c = Get()
	for WhiteSpace(c) {
		c = Get()
	}
	switch c {
	case EOF:
		token = EOF
	case '(', ')':
		token = c
	default:
		for i = 0; i < 100-1; {
			tokenbuf[i] = char(c)
			i++
			c = Get()
			if c == EOF {
				break
			}
			if WhiteSpace(c) || c == ')' {
				peekc = c
				break
			}
		}
		if i >= 100-1 {
			panic("token too long")
		}
		tokenbuf[i] = 0
		if '0' <= tokenbuf[0] && tokenbuf[0] <= '9' {
			token = '0'
		} else {
			token = 'A'
		}
	}
}

func Expect(c int) {
	if token != c {
		panic("parse error: expected " + string(c) + " got " + string(token))
	}
	NextToken()
}

func ParseList() *Slist {
	var slist, retval *Slist
	slist = new(Slist)
	slist.list.car = nil
	slist.list.cdr = nil
	slist.isatom = false
	slist.isstring = false
	retval = slist
	for {
		slist.list.car = Parse()
		if token == ')' {
			break
		}

		if token == EOF {
			break
		}

		slist.list.cdr = new(Slist)
		slist = slist.list.cdr
	}
	return retval
}

func atom(i int) *Slist {
	slist := new(Slist)
	if token == '0' {
		slist.atom.integer = i
		slist.isstring = false
	} else {
		slist.atom.string = new([100]char)
		var i int
		for i = 0; ; i++ {
			(*slist.atom.string)[i] = tokenbuf[i]
			if tokenbuf[i] == 0 {
				break
			}
		}
		slist.isstring = true
	}
	slist.isatom = true
	return slist
}

func atoi() int {
	var v int = 0
	for i := 0; '0' <= tokenbuf[i] && tokenbuf[i] <= '9'; i++ {
		v = v*10 + int(tokenbuf[i]-'0')
	}
	return v
}

func Parse() *Slist {
	var slist *Slist
	if token == EOF || token == ')' {
		return nil
	}

	if token == '(' {
		NextToken()
		slist = ParseList()
		Expect(')')
		return slist
	} else {
		switch token {
		case EOF:
			return nil
		case '0':
			slist = atom(atoi())
		case '"':
		case 'A':
			slist = atom(0)
		default:
			slist = nil
			fmt.Println("unknown token " + string(token))
		}
		NextToken()
		return slist
	}
}

func OpenFile() {
	inputindex = 0
	peekc = -1
	for i, c := range []byte(TestString) {
		input[i] = char(c)
	}
	NextToken()
}
