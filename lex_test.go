package lexpr

import (
	"context"
	"reflect"
	"testing"
)

func Test_lex_Parse(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []lexem
	}{
		{
			name: "math",
			args: args{
				input: "min(3, 2) * max(10, 20) == 40",
			},
			want: []lexem{
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
					Type:  op,
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
					Type:  op,
					Value: "==",
				}, {
					Type:  number,
					Value: "40",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLex()
			gotCh := l.parse(context.Background(), tt.args.input)
			got := []lexem{}
			for o := range gotCh {
				got = append(got, lexem{
					Type:  o.Type,
					Value: o.Value,
				})
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lex.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
