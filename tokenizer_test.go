package lexpr

import (
	"context"
	"reflect"
	"testing"
)

func TestLexpr_tokenize(t *testing.T) {
	type args struct {
		lexems []lexem
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			name: "math",
			args: args{
				lexems: []lexem{
					{
						Type:  word,
						Value: "min",
					}, {
						Type:  lp,
						Value: "(",
					}, {
						Type:  number,
						Value: "3",
					}, {
						Type:  sep,
						Value: ",",
					}, {
						Type:  number,
						Value: "2",
					}, {
						Type:  rp,
						Value: ")",
					}, {
						Type:  word,
						Value: "*",
					}, {
						Type:  word,
						Value: "max",
					}, {
						Type:  lp,
						Value: "(",
					}, {
						Type:  number,
						Value: "10",
					}, {
						Type:  sep,
						Value: ",",
					}, {
						Type:  number,
						Value: "20",
					}, {
						Type:  rp,
						Value: ")",
					}, {
						Type:  word,
						Value: "==",
					}, {
						Type:  number,
						Value: "40",
					},
				},
			},
			want: []Token{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexpr{
				operators: Operators,
				functions: Functions,
			}
			lexemsCh := make(chan lexem)
			go func() {
				defer close(lexemsCh)
				for _, l := range tt.args.lexems {
					lexemsCh <- l
				}
			}()
			gotCh := l.tokenize(context.Background(), lexemsCh)
			got := []Token{}
			for o := range gotCh {
				got = append(got, o)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lexpr.tokenize() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
