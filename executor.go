package lexpr

import (
	"context"
	"fmt"
	"strings"
)

func (l *Lexpr) execute(ctx context.Context, tokens <-chan Token) chan Result {
	out := make(chan Result)
	stack := TokenStack{}
	go func() {
		defer func() {
			for len(stack) > 0 {
				ret := stack.Pop()
				switch ret.typ {
				case str:
					out <- Result{Value: ret.value}
				case number:
					out <- Result{Value: ret.ivalue}
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
				case number:
					stack.Push(tkn)
				case str:
					stack.Push(Token{
						typ:   str,
						value: strings.Trim(tkn.value, `"`),
					})
				case funct:
					fn := l.functions[tkn.value]
					if err := fn(&stack); err != nil {
						out <- Result{Error: err}
						return
					}

				case op:
					op := l.operators[tkn.value]
					if err := op.handler(&stack); err != nil {
						out <- Result{Error: err}
						return
					}

				case word:
					variable, hasVariable := l.variables[strings.ToLower(tkn.value)]
					if !hasVariable {
						stack.Push(tkn)
						continue
					}
					vtkn, ok := TokenFromAny(variable)
					if !ok {
						out <- Result{Error: fmt.Errorf("invalid variable value: %+v", variable)}
						return
					}
					stack.Push(vtkn)
				case tokError:
					out <- Result{Error: fmt.Errorf(tkn.value)}
					return
				}
			}
		}
	}()
	return out
}
