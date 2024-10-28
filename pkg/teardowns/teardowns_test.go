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
		name string
		tr   *Funcs
		args args
	}{
		{
			name: "With valid function",
			tr:   New(),
			args: args{
				f: NewValidFunc(t),
			},
		},
		{
			name: "With nil teardowns",
			tr:   nil,
			args: args{
				f: NewValidFunc(t),
			},
		},
		{
			name: "With invalid function",
			tr:   New(),
			args: args{
				f: NewInvalidFunc(t),
			},
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
		name string
		tr   *Funcs
		args args
	}{
		{
			name: "With valid Funcs",
			tr:   New(),
			args: args{
				fs: []Func{
					NewValidFunc(t),
					NewValidFunc(t),
					NewValidFunc(t),
				},
			},
		},
		{
			name: "With one invalid Func",
			tr:   New(),
			args: args{
				fs: []Func{
					NewInvalidFunc(t),
					NewValidFunc(t),
					NewValidFunc(t),
				},
			},
		},
		{
			name: "With nil teardowns",
			tr:   nil,
			args: args{
				fs: []Func{
					NewValidFunc(t),
					NewValidFunc(t),
					NewValidFunc(t),
				},
			},
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
		tr   *Funcs
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
			tr: &Funcs{
				NewValidFunc(t),
				NewValidFunc(t),
			},
		},
		{
			name: "With 2 invalid teardowns",
			tr: &Funcs{
				NewInvalidFunc(t),
				NewInvalidFunc(t),
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
		name    string
		tr      *Funcs
		wantErr bool
		err     error
	}{
		{
			name: "With valid Funcs",
			tr: &Funcs{
				NewValidFunc(t),
				NewValidFunc(t),
				NewValidFunc(t),
			},
			wantErr: false,
		},
		{
			name: "With invalid Funcs",
			tr: &Funcs{
				NewInvalidFunc(t),
				NewInvalidFunc(t),
				NewInvalidFunc(t),
			},
			wantErr: true,
		},
		{
			name:    "With empty teardowns",
			tr:      New(),
			wantErr: false,
		},
		{
			name:    "With nil teardowns",
			tr:      nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tr.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Funcs.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFuncs_ToTeardownFunc(t *testing.T) {
	tests := []struct {
		name string
		tr   *Funcs
	}{
		{
			name: "With valid Funcs",
			tr: &Funcs{
				NewValidFunc(t),
				NewValidFunc(t),
				NewValidFunc(t),
			},
		},
		{
			name: "With invalid Funcs",
			tr: &Funcs{
				NewInvalidFunc(t),
				NewInvalidFunc(t),
				NewInvalidFunc(t),
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
				trc = *tt.tr
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

func NewValidFunc(t *testing.T) Func {
	return func() error {
		return nil
	}
}

func NewInvalidFunc(t *testing.T) Func {
	return func() error {
		return fmt.Errorf("invalid function")
	}
}
