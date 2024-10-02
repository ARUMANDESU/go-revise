package valueobject

import (
	"strings"
)

// Tags is a slice of strings that represents a list of tags.
type Tags []string

// NewTags creates a new Tags object with the given tags.
//
//	NOTE: If the tags array is empty or nil, it will return nil.
func NewTags(tags ...string) Tags {
	var t Tags
	t.AddMany(tags...)
	return t
}

// Contains checks if the given tag exists in the array.
func (t *Tags) Contains(want string) bool {
	if t == nil || len(*t) == 0 || want == "" {
		return false
	}

	for _, tag := range *t {
		if tag == want {
			return true
		}
	}

	return false
}

// Add adds a new tag to the array.
// If the tag already exists, it will not be added again.
// The tag will be trimmed before adding it to the array.
func (t *Tags) Add(tag string) {
	tag = strings.TrimSpace(tag)
	if t == nil || tag == "" {
		return
	}
	// if the tag already exists, do not add it again
	if t.Contains(tag) {
		return
	}

	*t = append(*t, tag)
}

// AddMany adds multiple tags to the array.
// If the tag already exists, it will not be added again.
func (t *Tags) AddMany(tags ...string) {
	for _, tag := range tags {
		t.Add(tag)
	}
}

// Remove removes a tag from the array.
// If the tag does not exist, it will not be removed.
func (t *Tags) Remove(tag string) {
	if t == nil || len(*t) == 0 {
		return
	}

	// trim the tag before comparing
	tag = strings.TrimSpace(tag)
	for i, v := range *t {
		if v == tag {
			// remove the tag from the array
			*t = append((*t)[:i], (*t)[i+1:]...)
			return
		}
	}
}

func (t *Tags) RemoveMany(tags ...string) {
	for _, tag := range tags {
		t.Remove(tag)
	}
}

// IsValid checks if the tags array is valid.
func (t *Tags) IsValid() bool {
	if t == nil || len(*t) == 0 {
		return false
	}
	return true
}

// TrimSpace trims spaces from each string in the array and returns a new array.
func (t *Tags) TrimSpace() Tags {
	if t == nil || len(*t) == 0 {
		return nil
	}

	// copy the array to avoid modifying the original
	arr := *t

	// trim spaces before converting to string
	for i, tag := range arr {
		arr[i] = strings.TrimSpace(tag)
	}

	return arr
}

// Unique returns a new array with unique tags.
func (t *Tags) Unique() Tags {
	if t == nil || len(*t) == 0 {
		return nil
	}

	// copy the array to avoid modifying the original
	arr := *t

	// create a map to store unique tags
	unique := make(map[string]struct{})

	// add tags to the map
	for _, tag := range arr {
		unique[tag] = struct{}{}
	}

	// convert the map to an array
	var tags Tags
	for tag := range unique {
		tags = append(tags, tag)
	}

	return tags
}

// Normalize trims spaces and removes duplicate tags from the array.
func (t *Tags) Normalize() {
	if t == nil || len(*t) == 0 {
		return
	}

	*t = t.TrimSpace()
	*t = t.Unique()
}
