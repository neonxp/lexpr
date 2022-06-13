package lexpr

import (
	"context"
	"fmt"
	"strconv"
)

func (l *Lexpr) tokenize(ctx context.Context, lexems <-chan lexem) <-chan Token {
	out := make(chan Token)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case lexem, ok := <-lexems:
				if !ok {
					return
				}
				switch {
				case lexem.Type == lp:
					out <- Token{
						typ: lp,
					}
				case lexem.Type == rp:
					out <- Token{
						typ: rp,
					}
				case lexem.Type == sep:
					out <- Token{
						typ: sep,
					}
				case lexem.Type == number:
					ivalue, _ := strconv.Atoi(lexem.Value)
					out <- Token{
						typ:    number,
						ivalue: ivalue,
					}
				case lexem.Type == str:
					out <- Token{
						typ:   str,
						value: lexem.Value,
					}
				case lexem.Type == op:
					o, isOp := l.operators[lexem.Value]
					if !isOp {
						out <- Token{
							typ:   tokError,
							value: fmt.Sprintf("unknown operator: %s", lexem.Value),
						}
						return
					}
					out <- Token{
						typ:       op,
						value:     lexem.Value,
						priority:  o.priority,
						leftAssoc: o.leftAssoc,
					}
				case lexem.Type == word:
					o, isOp := l.operators[lexem.Value]
					_, isFunc := l.functions[lexem.Value]
					switch {
					case isOp:
						out <- Token{
							typ:       op,
							value:     lexem.Value,
							priority:  o.priority,
							leftAssoc: o.leftAssoc,
						}
					case isFunc:
						out <- Token{
							typ:   funct,
							value: lexem.Value,
						}
					default:
						out <- Token{
							typ:   word,
							value: lexem.Value,
						}
					}
				case lexem.Type == tokError:
					out <- Token{
						typ:   tokError,
						value: lexem.Value,
					}
					return
				}
			}
		}
	}()
	return out
}
