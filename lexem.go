package lexpr

// lexem represents part of parsed string.
type lexem struct {
	Type  lexType // Type of Lexem.
	Value string  // Value of Lexem.
	Start int     // Start position at input string.
	End   int     // End position at input string.
}

// lexType represents type of current lexem.
type lexType int

// Some std lexem types
const (
	lexEOF lexType = iota
	tokError
	number
	str
	word
	op
	funct
	lp
	rp
	sep
)
