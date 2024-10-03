package valueobject

import (
	"testing"

	"github.com/clarify/subtest"
)

func TestValidateTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		tags        Tags
		errExpected bool
	}{
		{
			name: "With valid tags",
			tags: ExampleValidTags,
		},
		{
			name: "With empty tags",
			tags: Tags{},
		},
		{
			name: "With one tag",
			tags: NewTags("tag1"),
		},
		{
			name:        "With a lot of tags",
			tags:        ExampleInValidTags(),
			errExpected: true,
		},
		{
			name:        "With too long tag",
			tags:        NewTags(ExampleInvalidTag),
			errExpected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTags(tt.tags)
			if !tt.errExpected {
				t.Run("Expect no error", subtest.Value(err).NoError())
			} else {
				t.Run("Expect error", subtest.Value(err).Error())
			}
		})
	}
}
