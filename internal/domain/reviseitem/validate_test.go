package reviseitem

import (
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/clarify/subtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func TestValidateName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
	}{
		{
			name:     "With valid Name in English",
			username: validName(t, language.English),
		},
		{
			name:     "With valid Name in Kazakh",
			username: validName(t, language.Kazakh),
		},
		{
			name:     "With valid Name in Russian",
			username: validName(t, language.Russian),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.username)
			t.Run("Expected no error", subtest.Value(err).NoError())
		})
	}
}

func TestValidateName_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
	}{
		{
			name:     "With long Name",
			username: longName(t),
		},
		{
			name:     "With empty Name",
			username: "",
		},
		{
			name:     "With whitespace Name",
			username: " ",
		},
		{
			name:     "With tab Name",
			username: "\t",
		},
		{
			name:     "With newline Name",
			username: "\n",
		},
		{
			name:     "With carriage return Name",
			username: "\r",
		},
		{
			name:     "With multiple whitespace Name",
			username: "  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedErrorType := errs.ErrorTypeIncorrectInput
			err := validateName(tt.username)
			t.Run("Expected", func(t *testing.T) {
				require.Error(t, err, "error is expected")
				assert.IsType(t, &errs.Error{}, err, "expected error type")
				assert.True(t, errs.IsErrorType(err, expectedErrorType), "expected error type")
			})
		})
	}
}

func TestValidateDescription(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "With  description in English",
			description: validDescription(t, language.English),
		},
		{
			name:        "With  description in Kazakh",
			description: validDescription(t, language.Kazakh),
		},
		{
			name:        "With  description in Russian",
			description: validDescription(t, language.Russian),
		},
		{
			name:        "With empty description",
			description: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)
			t.Run("Expect no error", subtest.Value(err).NoError())
		})
	}
}

func TestValidateDescription_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "With long description",
			description: longDescription(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedErrorType := errs.ErrorTypeIncorrectInput
			err := validateDescription(tt.description)
			t.Run("Expected", func(t *testing.T) {
				require.Error(t, err, "error is expected")
				assert.IsType(t, &errs.Error{}, err, "expected error type")
				assert.True(t, errs.IsErrorType(err, expectedErrorType), "expected error type")
			})
		})
	}
}

func TestValidateNextRevisionAt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		nextRevisionAt time.Time
	}{
		{
			name:           "With valid nextRevisionAt",
			nextRevisionAt: validNextRevisionAt(t),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNextRevisionAt(tt.nextRevisionAt)
			t.Run("Expect no error", subtest.Value(err).NoError())
		})
	}
}

func TestValidateNextRevisionAt_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		nextRevisionAt time.Time
	}{
		{
			name:           "With zero nextRevisionAt",
			nextRevisionAt: time.Time{},
		},
		{
			name:           "With past nextRevisionAt",
			nextRevisionAt: time.Now().AddDate(0, 0, -1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedErrorType := errs.ErrorTypeUnknown
			err := validateNextRevisionAt(tt.nextRevisionAt)
			t.Run("Expected", func(t *testing.T) {
				require.Error(t, err, "error is expected")
				assert.IsType(t, &errs.Error{}, err, "expected error type")
				assert.True(t, errs.IsErrorType(err, expectedErrorType), "expected error type")
			})
		})
	}
}

func validName(t *testing.T, lang language.Tag) string {
	t.Helper()

	switch lang {
	case language.English:
		return "Go Chapter 1"
	case language.Kazakh:
		return "Go 1-тарау"
	case language.Russian:
		return "Глава 1 Go"
	default:
		return "Go Chapter 1" // by default return English
	}
}

func longName(t *testing.T) string {
	t.Helper()

	return strings.Repeat("a", maxNameLength+1)
}

func validDescription(t *testing.T, lang language.Tag) string {
	t.Helper()

	switch lang {
	case language.English:
		return "Complete Chapter 1 of Go Programming"
	case language.Kazakh:
		return "Go бағдарламасының 1-тарауын аяқтау"
	case language.Russian:
		return "Завершить Главу 1 по Go программированию"
	default:
		return "Complete Chapter 1 of Go Programming" // by default return English
	}
}

func validTags(t *testing.T) valueobject.Tags {
	t.Helper()

	return valueobject.NewTags("tag1", "tag2", "tag3")
}

func validNextRevisionAt(t *testing.T) time.Time {
	t.Helper()

	return time.Now().AddDate(0, 0, 1)
}

func longDescription(t *testing.T) string {
	t.Helper()

	return gofakeit.Sentence(maxDescriptionLength + 1)
}
