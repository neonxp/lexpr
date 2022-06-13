package lexpr

import (
	"context"
	"strings"
	"unicode/utf8"
)

// EOF const.
const EOF rune = -1

// lex holds current scanner state.
type lex struct {
	input  string     // Input string.
	start  int        // Start position of current lexem.
	pos    int        // Pos at input string.
	output chan lexem // Lexems channel.
	width  int        // Width of last rune.
}

// newLex returns new scanner for input string.
func newLex() *lex {
	return &lex{
		input:  "",
		start:  0,
		pos:    0,
		output: nil,
		width:  0,
	}
}

// parse input to lexems.
func (l *lex) parse(ctx context.Context, input string) <-chan lexem {
	l.input = input
	l.output = make(chan lexem)
	go func() {
		defer close(l.output)
		for {
			if ctx.Err() != nil {
				return
			}
			switch {
			case l.acceptWhile(" \n\t", false):
				l.ignore()
			case l.accept("("):
				l.emit(lp)
			case l.accept(")"):
				l.emit(rp)
			case l.accept(","):
				l.emit(sep)
			case scanNumber(l):
				l.emit(number)
			case scanOps(l):
				l.emit(op)
			case scanWord(l):
				l.emit(word)
			case scanQuotedString(l, `"`):
				l.emit(str)
			case l.peek() == EOF:
				return
			default:
				l.emit(tokError)
				return
			}
		}
	}()
	return l.output
}

// emit current lexem to output.
func (l *lex) emit(typ lexType) {
	l.output <- lexem{
		Type:  typ,
		Value: l.input[l.start:l.pos],
		Start: l.start,
		End:   l.pos,
	}
	l.start = l.pos
}

// next rune from input.
func (l *lex) next() (r rune) {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return EOF
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// back move position to previos rune.
func (l *lex) back() {
	l.pos -= l.width
}

// ignore previosly buffered text.
func (l *lex) ignore() {
	l.start = l.pos
	l.width = 0
}

// peek rune at current position without moving position.
func (l *lex) peek() (r rune) {
	r = l.next()
	l.back()
	return r
}

// accept any rune from valid string. Returns true if next rune was in valid string.
func (l *lex) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.back()
	return false
}

// acceptString returns true if given string was at position.
func (l *lex) acceptString(s string, caseInsentive bool) bool {
	input := l.input
	if caseInsentive {
		input = strings.ToLower(input)
		s = strings.ToLower(s)
	}
	if strings.HasPrefix(input, s) {
		l.width = 0
		l.pos += len(s)
		return true
	}
	return false
}

// acceptAnyOf substrings. Retuns true if any of substrings was found.
func (l *lex) acceptAnyOf(s []string, caseInsentive bool) bool {
	for _, substring := range s {
		if l.acceptString(substring, caseInsentive) {
			return true
		}
	}
	return false
}

// acceptWhile passing symbols from input while they at `valid` string.
func (l *lex) acceptWhile(valid string, ignoreEscaped bool) bool {
	start := l.pos
	for {
		ch := l.next()
		switch {
		case ch == EOF:
			return false
		case ch == '\\' && ignoreEscaped:
			l.next()
		case !strings.ContainsRune(valid, ch):
			l.back()
			return l.pos > start
		}
	}
}

// acceptWhileNot passing symbols from input while they NOT in `invalid` string.
func (l *lex) acceptWhileNot(invalid string, ignoreEscaped bool) bool {
	start := l.pos
	for {
		ch := l.next()
		switch {
		case ch == EOF:
			return false
		case ch == '\\' && ignoreEscaped:
			l.next()
		case strings.ContainsRune(invalid, ch):
			l.back()
			return l.pos > start
		}
	}
}

// atStart returns true if current lexem not empty
func (l *lex) atStart() bool {
	return l.pos == l.start
}
