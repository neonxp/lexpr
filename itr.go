package lexpr

import "context"

func infixToRpn(ctx context.Context, tokens <-chan Token) <-chan Token {
	out := make(chan Token)
	stack := TokenStack{}
	go func() {
		defer func() {
			if len(stack) > 0 {
				for {
					if stack.Head().typ == lp {
						out <- Token{
							typ:   tokError,
							value: "invalid brakets",
						}
						break
					}
					out <- stack.Pop()
					if len(stack) == 0 {
						break
					}
				}
			}
			close(out)
		}()
		for {
			select {
			case <-ctx.Done():
				return
			case tkn, ok := <-tokens:
				if !ok {
					return
				}
				switch tkn.typ {
				case number, word, str, tokError:
					out <- tkn
				case funct:
					stack.Push(tkn)
				case sep:
					for stack.Head().typ != lp {
						if len(stack) == 0 {
							out <- Token{
								typ:   tokError,
								value: "no arg separator or opening braket",
							}
							return
						}
						out <- stack.Pop()
					}
				case op:
					for len(stack) > 0 && (stack.Head().typ != op || (stack.Head().priority >= tkn.priority)) {
						out <- stack.Pop()
					}
					stack.Push(tkn)
				case lp:
					stack.Push(tkn)
				case rp:
					for stack.Head().typ != lp {
						if len(stack) == 0 {
							out <- Token{
								typ:   tokError,
								value: "no opening braket",
							}
							return
						}
						out <- stack.Pop()
					}
					stack.Pop()
					if stack.Head().typ == funct {
						out <- stack.Pop()
					}
				}
			}
		}
	}()
	return out
}
