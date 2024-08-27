package domain

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
)

func TestValidateFilterUserID(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"Valid UUID", gofakeit.UUID(), false},
		{"Invalid UUID", "invalid-uuid", true},
		{"Valid Int64", int64(1234567890), false},
		{"Invalid Type", 123.45, true},
		{"Invalid String", "not-a-uuid", true},
		{"Empty String", "", false}, // empty string is valid, because requirement is added additionaly
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFilterUserID(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{
			name: "valid",
			value: &Pagination{
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
		},
		{
			name:    "invalid Type",
			value:   "invalid",
			wantErr: true,
		},
		{
			name: "page less than 1",
			value: &Pagination{
				Page:     0,
				PageSize: 10,
			},
			wantErr: true,
		},
		{
			name: "page size less than 1",
			value: &Pagination{
				Page:     1,
				PageSize: 0,
			},
			wantErr: true,
		},
		{
			name:    "nil",
			value:   (*Pagination)(nil),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePagination(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateSort(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"valid", &Sort{Field: SortFieldDefault, Order: SortOrderDefault}, false},
		{"empty field", &Sort{Field: "", Order: SortOrderDefault}, true},
		{"empty order", &Sort{Field: SortFieldDefault, Order: ""}, true},
		{"invalid type", "invalid", true},
		{"nil", (*Sort)(nil), false},
		{"invalid field", &Sort{Field: "invalid", Order: "asc"}, true},
		{"invalid order", &Sort{Field: SortFieldDefault, Order: "invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSort(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
