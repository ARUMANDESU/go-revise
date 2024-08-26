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
	ID             uuid.UUID
	UserID         uuid.UUID
	Name           string
	Description    string
	Tags           StringArray
	Iteration      ReviseIteration
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastRevisedAt  time.Time
	NextRevisionAt time.Time
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

func (a StringArray) Value() driver.Value {
	// transform the array into a string: ["a","b","c"] -> "a,b,c"
	stringValue := strings.Join(a, ",")
	if len(stringValue) == 0 {
		return nil
	}

	return stringValue
}
