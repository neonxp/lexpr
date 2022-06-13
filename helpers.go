package lexpr

import "context"

func (l *Lexpr) OneResult(ctx context.Context, expression string) (any, error) {
	select {
	case r := <-l.Eval(ctx, expression):
		return r.Value, r.Error
	case <-ctx.Done():
		return nil, nil
	}
}
