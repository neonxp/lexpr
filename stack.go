package lexpr

type TokenStack []Token

func (s *TokenStack) Push(item Token) {
	*s = append(*s, item)
}

func (s *TokenStack) Pop() (item Token) {
	if len(*s) == 0 {
		return
	}

	*s, item = (*s)[:len(*s)-1], (*s)[len(*s)-1]
	return item
}

func (s *TokenStack) Head() (item Token) {
	if len(*s) == 0 {
		return
	}
	return (*s)[len(*s)-1]
}
