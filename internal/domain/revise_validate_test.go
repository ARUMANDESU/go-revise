package domain

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidateTags(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"ValidTags", []string{gofakeit.LetterN(ValidTagsMinLength + 1)}, false},
		{"TooShortTags", []string{gofakeit.LetterN(ValidTagsMinLength - 1)}, true},
		{"TooLongTags", []string{gofakeit.LetterN(ValidTagsMaxLength + 1)}, true},
		{"EmptyTags", []string{}, false},
		{"InvalidTags", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"ValidName", gofakeit.LetterN(ValidNameMinLength + 1), false},
		{"TooShortName", gofakeit.LetterN(ValidNameMinLength - 1), true},
		{"TooLongName", gofakeit.LetterN(ValidNameMaxLength + 1), true},
		{"EmptyName", "", false},
		{"InvalidName", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateName(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		value   any
		wantErr bool
	}{
		{"ValidDescription", gofakeit.LetterN(ValidDescriptionMinLength + 1), false},
		{"TooLongDescription", gofakeit.LetterN(ValidDescriptionMaxLength + 1), true},
		{"EmptyDescription", "", false},
		{"InvalidDescription", 123, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
