package reviseitem

import (
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/clarify/subtest"
	"golang.org/x/text/language"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
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
		want     error
	}{
		{
			name:     "With empty Name",
			username: "",
			want:     ErrInvalidArgument,
		},
		{
			name:     "With whitespace Name",
			username: " ",
			want:     ErrInvalidArgument,
		},
		{
			name:     "With tab Name",
			username: "\t",
			want:     ErrInvalidArgument,
		},
		{
			name:     "With newline Name",
			username: "\n",
			want:     ErrInvalidArgument,
		},
		{
			name:     "With carriage return Name",
			username: "\r",
			want:     ErrInvalidArgument,
		},
		{
			name:     "With multiple whitespace Name",
			username: "  ",
			want:     ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateName(tt.username)
			t.Run("Expected", subtest.Value(err).ErrorIs(tt.want))
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
		want        error
	}{
		{
			name:        "With long description",
			description: longDescription(t),
			want:        ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)
			t.Run("Expect error", subtest.Value(err).ErrorIs(tt.want))
		})
	}
}

func TestValidateTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags valueobject.StringArray
	}{
		{
			name: "With valid tags",
			tags: validTags(t),
		},
		{
			name: "With no tags",
			tags: valueobject.StringArray{},
		},
		{
			name: "With one tag",
			tags: valueobject.StringArray{"tag1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTags(tt.tags)
			t.Run("Expect no error", subtest.Value(err).NoError())
		})
	}
}

func TestValidateTags_Invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags valueobject.StringArray
		want error
	}{
		{
			name: "With too many tags",
			tags: moreTags(t),
			want: ErrInvalidArgument,
		},
		{
			name: "With too long tag",
			tags: valueobject.StringArray{longTag(t), "tag2"},
			want: ErrInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTags(tt.tags)
			t.Run("Expect error", subtest.Value(err).ErrorIs(tt.want))
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
		want           error
	}{
		{
			name:           "With zero nextRevisionAt",
			nextRevisionAt: time.Time{},
			want:           ErrInvalidArgument,
		},
		{
			name:           "With past nextRevisionAt",
			nextRevisionAt: time.Now().AddDate(0, 0, -1),
			want:           ErrNextRevisionAtInPast,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateNextRevisionAt(tt.nextRevisionAt)
			t.Run("Expect error", subtest.Value(err).ErrorIs(tt.want))
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

func validTags(t *testing.T) valueobject.StringArray {
	t.Helper()

	return valueobject.StringArray{"tag1", "tag2"}
}

func validNextRevisionAt(t *testing.T) time.Time {
	t.Helper()

	return time.Now().AddDate(0, 0, 1)
}

func longDescription(t *testing.T) string {
	t.Helper()

	return gofakeit.Sentence(maxDescriptionLength + 1)
}

func moreTags(t *testing.T) valueobject.StringArray {
	t.Helper()

	tags := valueobject.StringArray{}
	for i := 0; i < maxNumTags+1; i++ {
		tags = append(tags, gofakeit.Word())
	}

	return tags
}

func longTag(t *testing.T) string {
	t.Helper()

	return strings.Repeat("a", maxTagLength+1)
}
