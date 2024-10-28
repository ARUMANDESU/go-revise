package retry

import (
	"fmt"
	"testing"

	"github.com/clarify/subtest"
	"github.com/stretchr/testify/assert"
)

func TestWithMaxRetries(t *testing.T) {
	type args struct {
		maxRetries int
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		{
			name: "With 0 maxRetries",
			args: args{maxRetries: 0},
			want: Option{maxRetries: 0},
		},
		{
			name: "With 2 maxRetries",
			args: args{maxRetries: 2},
			want: Option{maxRetries: 2},
		},
		{
			name: "With -1 maxRetries",
			args: args{maxRetries: -1},
			want: Option{maxRetries: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WithMaxRetries(tt.args.maxRetries)
			t.Run(
				fmt.Sprintf("Expect %d maxRetries", tt.want.maxRetries),
				subtest.Value(got).DeepEqual(tt.want),
			)
		})
	}
}

func TestDo(t *testing.T) {
	type args struct {
		opts Option
	}
	tests := []struct {
		name            string
		args            args
		wantErr         bool
		expectedRetries int
	}{
		{
			name:            "With 0 maxRetries",
			args:            args{opts: WithMaxRetries(0)},
			wantErr:         false,
			expectedRetries: 2,
		},
		{
			name:            "With 5 maxRetries",
			args:            args{opts: WithMaxRetries(5)},
			wantErr:         false,
			expectedRetries: 5,
		},
		{
			name:            "With 11 maxRetries",
			args:            args{opts: WithMaxRetries(11)},
			wantErr:         true,
			expectedRetries: 11,
		},
		{
			name:            "With 100 maxRetries",
			args:            args{opts: WithMaxRetries(100)},
			wantErr:         true,
			expectedRetries: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var counter Counter = 0
			got := Do(func() error {
				return counter.Inc()
			}, tt.args.opts)

			if tt.wantErr {
				assert.Error(t, got)
			} else {
				assert.NoError(t, got)
			}
			assert.Equal(t, tt.expectedRetries, int(counter))
		})
	}
}

type Counter int

// Inc method increments counter, returns error if counter exceeds 10 or Counter is nil
func (c *Counter) Inc() error {
	if c == nil {
		return fmt.Errorf("counter is nil")
	}

	*c++

	if *c > 10 {
		return fmt.Errorf("exceeds 10")
	}

	return nil
}
