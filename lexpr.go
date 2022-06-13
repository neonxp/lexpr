package lexpr

import (
	"context"
	"strings"
)

type Lexpr struct {
	operators map[string]Operator
	functions map[string]func(ts *TokenStack) error
	variables map[string]any
}

func New(opts ...Opt) *Lexpr {
	l := &Lexpr{}
	for _, o := range opts {
		o(l)
	}
	return l
}

func (l *Lexpr) Eval(ctx context.Context, expression string) chan Result {
	lexer := newLex()
	lexems := lexer.parse(ctx, expression)
	tokens := l.tokenize(ctx, lexems)
	rpnTokens := infixToRpn(ctx, tokens)
	return l.execute(ctx, rpnTokens)
}

func (l *Lexpr) SetFunction(name string, fn func(ts *TokenStack) error) *Lexpr {
	l.functions[strings.ToLower(name)] = fn
	return l
}

func (l *Lexpr) SetOperator(name string, fn func(ts *TokenStack) error, priority int, leftAssoc bool) *Lexpr {
	l.operators[strings.ToLower(name)] = Operator{
		handler:   fn,
		priority:  priority,
		leftAssoc: leftAssoc,
	}
	return l
}

func (l *Lexpr) SetVariable(name string, value any) *Lexpr {
	l.variables[strings.ToLower(name)] = value
	return l
}

type Result struct {
	Value any
	Error error
}
