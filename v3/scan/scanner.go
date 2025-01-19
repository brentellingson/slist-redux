package scan

import (
	"bufio"
	"io"
)

type Scanner struct {
	reader  io.RuneReader
	Current rune
	EOF     bool
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(reader)}
}

// Advance reads the next rune from the scanner, and sets the Current field to that rune.
func (s *Scanner) Advance() {
	c, _, err := s.reader.ReadRune()
	if err == io.EOF {
		s.EOF = true
	} else if err != nil {
		panic(err)
	}
	s.Current = c
}
