package lexpr

type Token struct {
	typ       lexType
	value     string
	ivalue    int
	priority  int
	leftAssoc bool
}

func (t Token) Number() (int, bool) {
	return t.ivalue, t.typ == number
}

func (t Token) String() (string, bool) {
	return t.value, t.typ == str
}

func (t Token) Word() (string, bool) {
	return t.value, t.typ == word
}

func TokenFromAny(variable any) (Token, bool) {
	if s, ok := variable.(string); ok {
		return Token{
			typ:   str,
			value: s,
		}, true
	}
	if n, ok := variable.(int); ok {
		return Token{
			typ:    number,
			ivalue: n,
		}, true
	}
	if n, ok := variable.(float64); ok {
		return Token{
			typ:    number,
			ivalue: int(n),
		}, true
	}
	if n, ok := variable.(float32); ok {
		return Token{
			typ:    number,
			ivalue: int(n),
		}, true
	}
	if b, ok := variable.(bool); ok {
		n := 0
		if b {
			n = 1
		}
		return Token{
			typ:    number,
			ivalue: n,
		}, true
	}
	return Token{}, false
}

func TokenFromWord(wordName string) Token {
	return Token{
		typ:   word,
		value: wordName,
	}
}

func TokenFromString(s string) Token {
	return Token{
		typ:   str,
		value: s,
	}
}

func TokenFromInt(n int) Token {
	return Token{
		typ:    number,
		ivalue: n,
	}
}
