package valueobject

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// StringArray is a custom type to handle string arrays in the database, since sqlite does not support arrays.
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	// Scan a database value into a string array: "a,b,c" -> ["a","b","c"]
	const op = "StringArray.Scan"

	if value == nil { // case when value from the db was NULL
		*a = nil
		return nil
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

// Value converts the string array into a string so that it can be saved in the database
//
//   - trims spaces before converting to string
//
//   - returns nil if the array is empty
//
//     Note: dos not modify the original array
func (a *StringArray) Value() driver.Value {
	// transform the array into a string: ["a","b","c"] -> "a,b,c"
	if a == nil || len(*a) == 0 {
		return nil
	}

	// copy the array to avoid modifying the original
	arr := *a

	// trim spaces before converting to string
	for i, tag := range arr {
		arr[i] = strings.TrimSpace(tag)
	}

	// convert to string
	stringValue := strings.Join(arr, ",")
	if len(stringValue) == 0 {
		return nil
	}

	return stringValue
}
