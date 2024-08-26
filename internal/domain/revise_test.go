package domain

import (
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringArray_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    StringArray
		wantErr bool
	}{
		{
			name:    "nil value",
			input:   nil,
			want:    nil,
			wantErr: false,
		},
		{
			name:    "valid string",
			input:   "tag1,tag2,tag3",
			want:    StringArray{"tag1", "tag2", "tag3"},
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   123,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var a StringArray
			err := a.Scan(tt.input)
			require.Equal(t, tt.wantErr, err != nil)

			assert.Equal(t, tt.want, a)
		})
	}
}

func TestStringArray_Value(t *testing.T) {
	tests := []struct {
		name string
		a    StringArray
		want driver.Value
	}{
		{
			name: "multiple tags",
			a:    StringArray{"tag1", "tag2", "tag3"},
			want: "tag1,tag2,tag3",
		},
		{
			name: "single tag",
			a:    StringArray{"tag1"},
			want: "tag1",
		},
		{
			name: "empty array",
			a:    StringArray{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.a.Value())
		})
	}
}
