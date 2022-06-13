package lexpr

import (
	"context"
	"reflect"
	"testing"
)

func TestLexpr_Eval(t *testing.T) {
	type fields struct {
		operators map[string]Operator
		functions map[string]func(ts *TokenStack) error
		variables map[string]any
	}
	type args struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "simple math",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{},
			},
			args:    args{expression: "2 + 2 * 2"},
			want:    6,
			wantErr: false,
		},
		{
			name: "complex equal",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{},
			},
			args:    args{expression: "min(3, 2) * max(10, 20) == 40"},
			want:    1,
			wantErr: false,
		},
		{
			name: "complex neql",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{},
			},
			args:    args{expression: "min(3, 2) * max(10, 20) != 40"},
			want:    0,
			wantErr: false,
		},
		{
			name: "variables",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{
					"svar": "test",
					"ivar": int(123),
					"fvar": 321.0,
				},
			},
			args: args{
				expression: "len(svar) + ivar + fvar",
			},
			want:    448,
			wantErr: false,
		},
		{
			name: "invalid1",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{},
			},
			args:    args{expression: ")("},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid2",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{},
			},
			args:    args{expression: "var1 + var2"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid3",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{},
			},
			args:    args{expression: "3 @ 4"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "dot notation",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{
					"j": `{ "one" : { "four": {"five": "six"} }, "two": "three" }`,
				},
			},
			args: args{
				expression: `j.one.four.five`,
			},
			want:    `six`,
			wantErr: false,
		},
		{
			name: "dot notation with arrays",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{
					"j": `{ "one" : { "four": ["five", "six", "seven"] }, "two": "three" }`,
				},
			},
			args: args{
				expression: `j.one.four.1`,
			},
			want:    `six`,
			wantErr: false,
		},
		{
			name: "dot notation with arrays and variables",
			fields: fields{
				operators: Operators,
				functions: Functions,
				variables: map[string]any{
					"j":    `{ "one" : { "four": ["five", "six", "seven"] }, "two": "three" }`,
					"key1": "one",
					"key2": 1,
				},
			},
			args: args{
				expression: `j.key1.four.key2`,
			},
			want:    `six`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexpr{
				operators: tt.fields.operators,
				functions: tt.fields.functions,
				variables: tt.fields.variables,
			}
			gotCh := l.Eval(context.Background(), tt.args.expression)
			res := <-gotCh
			got := res.Value
			err := res.Error
			if (err != nil) != tt.wantErr {
				t.Errorf("Lexpr.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexpr.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
