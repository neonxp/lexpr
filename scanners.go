package lexpr

import (
	"strings"
)

const (
	digits = "0123456789"
	alpha  = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	chars  = "+-*/=<>@&|:!."
)

// scanNumber simplest scanner that accepts decimal int and float.
func scanNumber(l *lex) bool {
	l.acceptWhile(digits, false)
	if l.atStart() {
		// not found any digit
		return false
	}
	l.accept(".")
	l.acceptWhile(digits, false)
	return !l.atStart()
}

// scanWord returns true if next input token contains alphanum sequence that not starts from digit and not contains.
// spaces or special characters.
func scanWord(l *lex) bool {
	if !l.accept(alpha) {
		return false
	}
	l.acceptWhile(alpha+digits, false)
	return true
}

func scanOps(l *lex) bool {
	return l.acceptWhile(chars, false)
}

// scanQuotedString returns true if next input tokens is quoted string. Can be used with any type of quotes.
func scanQuotedString(l *lex, quote string) bool {
	start := l.pos
	if !strings.ContainsRune(quote, l.next()) {
		l.pos = start
		return false
	}
	if l.acceptWhileNot(quote, true) {
		l.next()
	}
	return !l.atStart()
}
