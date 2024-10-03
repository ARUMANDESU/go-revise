// This folder user for testing purpose only

package valueobject

import (
	"github.com/brianvoe/gofakeit/v7"
)

// ---- Tags ----
var (
	// ExampleValidTags is an example of valid tags. It is used in tests.
	ExampleValidTags = NewTags("tag1", "tag2", "tag3")
	// ExampleInValidTags has more than the maximum number of tags. It is used in tests.
	ExampleInValidTags = func() Tags {
		tags := make([]string, maxNumTags+1)
		for i := 0; i < maxNumTags+1; i++ {
			tags[i] = gofakeit.LetterN(uint(maxTagLength))
		}
		return NewTags(tags...)
	}
	// ExampleInvalidTag is tag with more than the maximum number of characters. It is used in tests.
	ExampleInvalidTag = gofakeit.LetterN(uint(maxTagLength + 1))
)
