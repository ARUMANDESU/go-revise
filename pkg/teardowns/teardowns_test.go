package teardowns

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	got := New()

	require.NotNil(t, got)
}

func TestFuncs_Append(t *testing.T) {
	type args struct {
		f Func
	}
	tests := []struct {
		name        string
		tr          Funcs
		args        args
		expectedLen int
	}{
		{
			name: "With valid function",
			tr:   New(),
			args: args{
				f: newValidFunc(t),
			},
			expectedLen: 1,
		},
		{
			name: "With nil teardowns",
			tr:   nil,
			args: args{
				f: newValidFunc(t),
			},
			expectedLen: 1,
		},
		{
			name: "With invalid function",
			tr:   New(),
			args: args{
				f: newInvalidFunc(t),
			},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.Append(tt.args.f)
		})
	}
}

func TestFuncs_AppendMany(t *testing.T) {
	type args struct {
		fs []Func
	}
	tests := []struct {
		name        string
		tr          Funcs
		args        args
		expectedLen int
	}{
		{
			name: "With valid Funcs",
			tr:   New(),
			args: args{
				fs: []Func{
					newValidFunc(t),
					newValidFunc(t),
					newValidFunc(t),
				},
			},
			expectedLen: 3,
		},
		{
			name: "With one invalid Func",
			tr:   New(),
			args: args{
				fs: []Func{
					newInvalidFunc(t),
					newValidFunc(t),
					newValidFunc(t),
				},
			},
			expectedLen: 3,
		},
		{
			name: "With nil teardowns",
			tr:   nil,
			args: args{
				fs: []Func{
					newValidFunc(t),
					newValidFunc(t),
					newValidFunc(t),
				},
			},
			expectedLen: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.AppendMany(tt.args.fs...)
		})
	}
}

func TestFuncs_Clear(t *testing.T) {
	tests := []struct {
		name string
		tr   Funcs
	}{
		{
			name: "With empty teardowns",
			tr:   New(),
		},
		{
			name: "With nil teardowns",
			tr:   nil,
		},
		{
			name: "With 2 teardowns",
			tr: Funcs{
				newValidFunc(t),
				newValidFunc(t),
			},
		},
		{
			name: "With 2 invalid teardowns",
			tr: Funcs{
				newInvalidFunc(t),
				newInvalidFunc(t),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.tr.Clear()
		})
	}
}

func TestFuncs_Execute(t *testing.T) {
	tests := []struct {
		name                  string
		initFunc              func(t *testing.T) (Funcs, *testFunc)
		expectedNumberOfFuncs int
		wantErr               bool
		err                   error
	}{
		{
			name: "With valid Funcs",
			initFunc: func(t *testing.T) (Funcs, *testFunc) {
				t.Helper()

				tf := &testFunc{}
				tr := Funcs{
					tf.Func(true),
					tf.Func(true),
					tf.Func(true),
				}
				return tr, tf
			},
			expectedNumberOfFuncs: 3,
			wantErr:               false,
		},
		{
			name: "With invalid Funcs",
			initFunc: func(t *testing.T) (Funcs, *testFunc) {
				t.Helper()

				tf := &testFunc{}
				tr := Funcs{
					tf.Func(false),
					tf.Func(false),
					tf.Func(false),
				}
				return tr, tf
			},
			expectedNumberOfFuncs: 3,
			wantErr:               true,
		},
		{
			name: "With empty teardowns",
			initFunc: func(t *testing.T) (Funcs, *testFunc) {
				t.Helper()
				tr := New()
				return tr, &testFunc{}
			},
			expectedNumberOfFuncs: 0,
			wantErr:               false,
		},
		{
			name: "With nil teardowns",
			initFunc: func(t *testing.T) (Funcs, *testFunc) {
				t.Helper()
				return nil, &testFunc{}
			},
			expectedNumberOfFuncs: 0,
			wantErr:               false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, tf := tt.initFunc(t)
			if err := tr.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Funcs.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.expectedNumberOfFuncs, tf.counter)
		})
	}
}
func TestFuncs_ToTeardownFunc(t *testing.T) {
	tests := []struct {
		name string
		tr   Funcs
	}{
		{
			name: "With valid Funcs",
			tr: Funcs{
				newValidFunc(t),
				newValidFunc(t),
				newValidFunc(t),
			},
		},
		{
			name: "With invalid Funcs",
			tr: Funcs{
				newInvalidFunc(t),
				newInvalidFunc(t),
				newInvalidFunc(t),
			},
		},
		{
			name: "With empty teardowns",
			tr:   New(),
		},
		{
			name: "With nil teardowns",
			tr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// copy teardowns so Execute dose not take effect on ToTeardownFunc
			var trc Funcs
			if tt.tr != nil {
				trc = tt.tr
			}

			want := trc.Execute()
			got := tt.tr.ToTeardownFunc()()
			if want != nil {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

type testFunc struct {
	counter int
}

func (tf *testFunc) Func(isValid bool) Func {
	return func() error {
		tf.counter++
		if isValid {
			return nil
		}
		return fmt.Errorf("invalid function")
	}
}

func newValidFunc(t *testing.T) Func {
	t.Helper()
	return func() error {
		return nil
	}
}

func newInvalidFunc(t *testing.T) Func {
	t.Helper()
	return func() error {
		return fmt.Errorf("invalid function")
	}
}
