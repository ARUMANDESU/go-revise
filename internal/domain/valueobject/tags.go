package valueobject

import (
	"strings"
)

// Tags represents a list of tags encapsulated in a struct.
type Tags struct {
	tags []string
}

// NewTags creates a new Tags object with the given tags.
//
//	NOTE: If the tags array is empty or nil, it will return nil.
func NewTags(tags ...string) Tags {
	t := Tags{}
	t.AddMany(tags...)
	return t
}

// IsTagsEqual checks if two tags are equal.
func IsTagsEqual(a, b *Tags) bool {
	if a == nil || b == nil || len(a.tags) != len(b.tags) {
		return false
	}

	for _, tag := range a.tags {
		// Check if the tags are the same, but the order can be different.
		if !b.Contains(tag) {
			return false
		}
	}

	return true
}

func (t *Tags) StringArray() []string {
	return t.tags
}

// Contains checks if the given tag exists in the list of tags.
func (t *Tags) Contains(want string) bool {
	if t == nil || len(t.tags) == 0 || want == "" {
		return false
	}

	for _, tag := range t.tags {
		if tag == want {
			return true
		}
	}

	return false
}

// Add adds a new tag to the list.
// If the tag already exists, it will not be added again.
// The tag will be trimmed before adding it to the list.
func (t *Tags) Add(tag string) {
	tag = strings.TrimSpace(tag)
	if t == nil || tag == "" {
		return
	}
	// if the tag already exists, do not add it again
	if t.Contains(tag) {
		return
	}

	t.tags = append(t.tags, tag)
}

// AddMany adds multiple tags to the list.
// If a tag already exists, it will not be added again.
func (t *Tags) AddMany(tags ...string) {
	for _, tag := range tags {
		t.Add(tag)
	}
}

// AddTags adds multiple tags to the list.
func (t *Tags) AddTags(tags Tags) {
	t.AddMany(tags.tags...)
}

// Remove removes a tag from the list.
// If the tag does not exist, it will not be removed.
func (t *Tags) Remove(tag string) {
	if t == nil || len(t.tags) == 0 {
		return
	}

	// trim the tag before comparing
	tag = strings.TrimSpace(tag)
	for i, v := range t.tags {
		if v == tag {
			// remove the tag from the list
			t.tags = append(t.tags[:i], t.tags[i+1:]...)
			return
		}
	}
}

// RemoveMany removes multiple tags from the list.
func (t *Tags) RemoveMany(tags ...string) {
	for _, tag := range tags {
		t.Remove(tag)
	}
}

// RemoveTags removes multiple tags from the list.
func (t *Tags) RemoveTags(tags Tags) {
	t.RemoveMany(tags.tags...)
}

func (t *Tags) IsEmpty() bool {
	if t == nil || len(t.tags) == 0 {
		return false
	}
	return true
}

// TrimSpace trims spaces from each string in the list and returns a new list.
func (t *Tags) TrimSpace() Tags {
	if t == nil || len(t.tags) == 0 {
		return Tags{}
	}

	// copy the list to avoid modifying the original
	arr := t.tags

	// trim spaces before converting to string
	for i, tag := range arr {
		arr[i] = strings.TrimSpace(tag)
	}

	return Tags{tags: arr}
}

// Unique returns a new list with unique tags.
func (t *Tags) Unique() Tags {
	if t == nil || len(t.tags) == 0 {
		return Tags{}
	}

	// create a map to store unique tags
	unique := make(map[string]struct{})

	// add tags to the map
	for _, tag := range t.tags {
		unique[tag] = struct{}{}
	}

	// convert the map to a list
	var tags []string
	for tag := range unique {
		tags = append(tags, tag)
	}

	return Tags{tags: tags}
}

// Normalize trims spaces and removes duplicate tags from the list.
func (t *Tags) Normalize() {
	if t == nil || len(t.tags) == 0 {
		return
	}

	t.tags = t.TrimSpace().tags
	t.tags = t.Unique().tags
}
