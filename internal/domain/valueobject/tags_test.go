package valueobject

import (
	"fmt"
	"testing"

	"github.com/clarify/subtest"
)

func TestNewTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tags []string
		want Tags
	}{
		{
			name: "With tags",
			tags: []string{"tag1", "tag2"},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With empty tags",
			tags: []string{},
			want: nil,
		},
		{
			name: "With nil tags",
			tags: nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTags(tt.tags...)
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(got, tt.want))
		})
	}
}

func TestTags_Contains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		t        Tags
		want     string
		expected bool
	}{
		{
			name:     "With existing tag",
			t:        Tags{"tag1", "tag2"},
			want:     "tag1",
			expected: true,
		},
		{
			name:     "With non-existing tag",
			t:        Tags{"tag1", "tag2"},
			want:     "tag3",
			expected: false,
		},
		{
			name:     "With nil tags",
			t:        nil,
			want:     "tag1",
			expected: false,
		},
		{
			name:     "With empty tags",
			t:        Tags{},
			want:     "tag1",
			expected: false,
		},
		{
			name:     "With empty tag",
			t:        Tags{"tag1", "tag2"},
			want:     "",
			expected: false,
		},
		{
			name:     "With empty tag and empty tags",
			t:        Tags{},
			want:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.Contains(tt.want)
			t.Run(fmt.Sprintf("Expected %v", tt.expected), subtest.Value(got).CompareEqual(tt.expected))
		})
	}
}

func TestTags_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		tag  string
		want Tags
	}{
		{
			name: "With existing tag",
			t:    Tags{"tag1", "tag2"},
			tag:  "tag1",
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With new tag",
			t:    Tags{"tag1", "tag2"},
			tag:  "tag3",
			want: Tags{"tag1", "tag2", "tag3"},
		},
		{
			name: "With empty tag",
			t:    Tags{"tag1", "tag2"},
			tag:  "",
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With empty tag and empty tags",
			t:    Tags{},
			tag:  "",
			want: Tags{},
		},
		{
			name: "With nil tags",
			t:    nil,
			tag:  "tag1",
			want: Tags{"tag1"},
		},
		{
			name: "With empty tags",
			t:    Tags{},
			tag:  "tag1",
			want: Tags{"tag1"},
		},
		{
			name: "With empty tag and empty tags",
			t:    Tags{},
			tag:  "",
			want: Tags{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Add(tt.tag)
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(tt.t, tt.want))
		})
	}
}

func TestTags_AddMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		tags []string
		want Tags
	}{
		{
			name: "With existing tags",
			t:    Tags{"tag1", "tag2"},
			tags: []string{"tag1", "tag2"},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With new tags",
			t:    Tags{"tag1", "tag2"},
			tags: []string{"tag3", "tag4"},
			want: Tags{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name: "With empty tags",
			t:    Tags{"tag1", "tag2"},
			tags: []string{},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With nil tags",
			t:    nil,
			tags: []string{"tag1", "tag2"},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With empty tags",
			t:    Tags{},
			tags: []string{"tag1", "tag2"},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With empty tags and empty tags",
			t:    Tags{},
			tags: []string{},
			want: Tags{},
		},
		{
			name: "With existing and new tags",
			t:    Tags{"tag1", "tag2"},
			tags: []string{"tag2", "tag3"},
			want: Tags{"tag1", "tag2", "tag3"},
		},
		{
			name: "With same tags",
			t:    Tags{"tag1", "tag2"},
			tags: []string{"tag1", "tag2", "tag2", "tag1"},
			want: Tags{"tag1", "tag2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.AddMany(tt.tags...)
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(tt.t, tt.want))
		})
	}
}

func TestTags_Remove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		tag  string
		want Tags
	}{
		{
			name: "With existing tag",
			t:    Tags{"tag1", "tag2"},
			tag:  "tag1",
			want: Tags{"tag2"},
		},
		{
			name: "With non-existing tag",
			t:    Tags{"tag1", "tag2"},
			tag:  "tag3",
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With empty tag",
			t:    Tags{"tag1", "tag2"},
			tag:  "",
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With nil tags",
			t:    nil,
			tag:  "tag1",
			want: nil,
		},
		{
			name: "With empty tags",
			t:    Tags{},
			tag:  "tag1",
			want: Tags{},
		},
		{
			name: "With empty tag and empty tags",
			t:    Tags{},
			tag:  "",
			want: Tags{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Remove(tt.tag)
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(tt.t, tt.want))
		})
	}
}

func TestTags_RemoveMany(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		tags []string
		want Tags
	}{
		{
			name: "With existing tags",
			t:    Tags{"tag1", "tag2", "tag3", "tag4"},
			tags: []string{"tag1", "tag2"},
			want: Tags{"tag3", "tag4"},
		},
		{
			name: "With non-existing tags",
			t:    Tags{"tag1", "tag2", "tag3", "tag4"},
			tags: []string{"tag5", "tag6"},
			want: Tags{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name: "With empty tags",
			t:    Tags{"tag1", "tag2", "tag3", "tag4"},
			tags: []string{},
			want: Tags{"tag1", "tag2", "tag3", "tag4"},
		},
		{
			name: "With nil tags",
			t:    nil,
			tags: []string{"tag1", "tag2"},
			want: nil,
		},
		{
			name: "With empty tags",
			t:    Tags{},
			tags: []string{"tag1", "tag2"},
			want: Tags{},
		},
		{
			name: "With empty tags and empty tags",
			t:    Tags{},
			tags: []string{},
			want: Tags{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.t.RemoveMany(tt.tags...)
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(tt.t, tt.want))
		})
	}
}

func TestTags_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		t        Tags
		expected bool
	}{
		{
			name:     "With existing tags",
			t:        Tags{"tag1", "tag2"},
			expected: true,
		},
		{
			name:     "With nil tags",
			t:        nil,
			expected: false,
		},
		{
			name:     "With empty tags",
			t:        Tags{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.IsValid()
			t.Run(fmt.Sprintf("Expected %v", tt.expected), subtest.Value(got).CompareEqual(tt.expected))
		})
	}
}

func TestTags_TrimSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		want Tags
	}{
		{
			name: "With existing tags",
			t:    Tags{" tag1", "tag2 ", " tag3 "},
			want: Tags{"tag1", "tag2", "tag3"},
		},
		{
			name: "With nil tags",
			t:    nil,
			want: nil,
		},
		{
			name: "With empty tags",
			t:    Tags{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.TrimSpace()
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(got, tt.want))
		})
	}
}

func TestTags_Unique(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		want Tags
	}{
		{
			name: "With existing tags",
			t:    Tags{"tag1", "tag2", "tag1", "tag2"},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With nil tags",
			t:    nil,
			want: nil,
		},
		{
			name: "With empty tags",
			t:    Tags{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.t.Unique()
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(got, tt.want))
		})
	}
}

func TestTags_Normalize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		t    Tags
		want Tags
	}{
		{
			name: "With existing tags",
			t:    Tags{"tag1", "tag2 ", " tag1", " tag2 "},
			want: Tags{"tag1", "tag2"},
		},
		{
			name: "With nil tags",
			t:    nil,
			want: nil,
		},
		{
			name: "With empty tags",
			t:    Tags{},
			want: Tags{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.t.Normalize()
			t.Run(fmt.Sprintf("Expected %v", tt.want), assertTags(tt.t, tt.want))
		})
	}
}

func assertTags(got, want Tags) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		if len(got) != len(want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
		// Check if the tags are the same, but the order can be different.
		for _, tag := range got {
			if !want.Contains(tag) {
				t.Errorf("got: %v, want: %v", got, want)
			}
		}
	}
}
