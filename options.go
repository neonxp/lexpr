package lexpr

type Opt func(*Lexpr)

func WithOperators(operators map[string]Operator) Opt {
	return func(l *Lexpr) {
		l.operators = operators
	}
}

func WithFunctions(functions map[string]func(ts *TokenStack) error) Opt {
	return func(l *Lexpr) {
		l.functions = functions
	}
}

func WithValues(variables map[string]any) Opt {
	return func(l *Lexpr) {
		l.variables = variables
	}
}

func WithDefaults() Opt {
	return func(l *Lexpr) {
		l.operators = Operators
		l.functions = Functions
		l.variables = map[string]any{}
	}
}
