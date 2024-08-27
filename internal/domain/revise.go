package domain

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// ReviseItem is represents revision entity
type ReviseItem struct {
	// I think optemistic locking is not needed here
	// because we are not updating the same record concurrently
	// so I am not adding version field here

	// Immutable fields
	ID     uuid.UUID
	UserID uuid.UUID

	// Updatable fields
	Name        string
	Description string
	Tags        StringArray

	// Revision specific fields
	Iteration      ReviseIteration
	LastRevisedAt  time.Time
	NextRevisionAt time.Time

	// Timestamps
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AbleToUpdate checks if the user is able to update the revise item
func (r ReviseItem) AbleToUpdate(id string) bool {
	return r.UserID.String() == id
}

// PartialUpdate updates only the fields that are specified in the dto.UpdateFields
// and returns the updated ReviseItem
//   - dto.Name is trimmed before updating
//   - dto.Description is trimmed before updating
//   - dto.Tags is not trimmed here, because it is handled in the StringArray.Value method, before saving to the database
func (r ReviseItem) PartialUpdate(dto UpdateReviseItemDTO) ReviseItem {
	for _, field := range dto.UpdateFields {
		switch field {
		case "name":
			r.Name = strings.TrimSpace(dto.Name)
		case "description":
			r.Description = strings.TrimSpace(dto.Description)
		case "tags":
			r.Tags = dto.Tags // no need to trim spaces here, because it is handled in the StringArray.Value method
		}
	}

	// update the updated at timestamp if any of the fields are updated
	if dto.UpdateFields != nil || len(dto.UpdateFields) > 0 {
		r.UpdatedAt = time.Now()
	}

	return r
}

// StringArray is a custom type to handle string array in the database,
// because sqlite does not support array types
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	// Scan a database value into a string array: "a,b,c" -> ["a","b","c"]
	const op = "domain.StringArray.Scan"

	if value == nil {
		return nil // case when value from the db was NULL
	}

	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("%s: failed to convert value to string", op)
	}
	tags := strings.Split(s, ",")
	if len(tags) == 0 || tags[0] == "" {
		return nil
	}
	*a = tags
	return nil
}

// Value converts the string array into a string
// so that it can be saved in the database
//   - trims spaces before converting to string
//
// Use this method before saving to the database
func (a StringArray) Value() driver.Value {
	// transform the array into a string: ["a","b","c"] -> "a,b,c"

	// trim spaces before converting to string
	for i, tag := range a {
		a[i] = strings.TrimSpace(tag)
	}

	stringValue := strings.Join(a, ",")
	if len(stringValue) == 0 {
		return nil
	}

	return stringValue
}
