package lexpr

import (
	"context"
	"reflect"
	"testing"
)

func Test_infixToRpn(t *testing.T) {
	type args struct {
		in []Token
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			name: "simple",
			args: args{
				in: []Token{
					{
						typ:   funct,
						value: "min",
					},
					{
						typ: lp,
					},
					{
						typ:    number,
						ivalue: 3,
					},
					{
						typ: sep,
					},
					{
						typ:    number,
						ivalue: 2,
					},
					{
						typ: rp,
					},
					{
						typ:       op,
						value:     "*",
						ivalue:    0,
						priority:  120,
						leftAssoc: false,
					},
					{
						typ:   funct,
						value: "max",
					},
					{
						typ: lp,
					},
					{
						typ:    number,
						ivalue: 10,
					},
					{
						typ: sep,
					},
					{
						typ:    number,
						ivalue: 20,
					},
					{
						typ: rp,
					},
					{
						typ:       op,
						value:     "==",
						ivalue:    0,
						priority:  20,
						leftAssoc: false,
					},
					{
						typ:    number,
						ivalue: 40,
					},
				},
			},
			want: []Token{
				{
					typ:    number,
					ivalue: 3,
				},
				{
					typ:    number,
					ivalue: 2,
				},
				{
					typ:   funct,
					value: "min",
				},
				{
					typ:    number,
					ivalue: 10,
				},
				{
					typ:    number,
					ivalue: 20,
				},
				{
					typ:   funct,
					value: "max",
				},
				{
					typ:       op,
					value:     "*",
					ivalue:    0,
					priority:  120,
					leftAssoc: false,
				},
				{
					typ:    number,
					ivalue: 40,
				},
				{
					typ:       op,
					value:     "==",
					ivalue:    0,
					priority:  20,
					leftAssoc: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan Token)
			go func() {
				defer close(inCh)
				for _, tk := range tt.args.in {
					inCh <- tk
				}
			}()
			gotCh := infixToRpn(context.Background(), inCh)
			got := []Token{}
			for o := range gotCh {
				got = append(got, o)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("infixToRpn() = %v, want %v", got, tt.want)
			}
		})
	}
}
