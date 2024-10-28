package pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "With int returns pointer to int",
			args: args{v: 42},
		},
		{
			name: "With string returns pointer to string",
			args: args{v: "hello"},
		},
		{
			name: "With struct returns pointer to struct",
			args: args{v: struct{ Field string }{Field: "value"}},
		},
		{
			name: "With nil interface returns pointer to nil interface",
			args: args{v: (interface{})(nil)},
		},
		{
			name: "With slice returns pointer to slice",
			args: args{v: []int{1, 2, 3}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.v)
			assert.Equal(t, tt.args.v, *got)
		})
	}
}
