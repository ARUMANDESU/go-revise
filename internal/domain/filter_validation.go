package domain

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func ValidateFilterUserID(value any) error {
	// user id can be either string(uuid) or int64(telegram user id)
	uid, ok := value.(string)
	if !ok {
		if _, ok := value.(int64); !ok {
			return errors.New("user id must be either string(uuid) or int64(telegram user id)")
		}

		return nil
	}

	return validation.Validate(uid, is.UUID) // if string is provided, it must be a valid uuid
}

func ValidatePagination(value any) error {
	pagination, ok := value.(*Pagination)
	if !ok {
		return errors.New("invalid type, must be pointer to Pagination")
	}
	// pagination is not required, so it can be nil
	if pagination == nil {
		return nil
	}

	if pagination.Page <= 0 {
		return errors.New("page must be greater than 0")
	}
	if pagination.PageSize <= 0 {
		return errors.New("page size must be greater than 0")
	}

	return nil
}

func ValidateSort(value any) error {
	sort, ok := value.(*Sort)
	if !ok {
		return errors.New("invalid sort type, must be pointer to Sort")
	}

	// sort is not required, so it can be nil
	if sort == nil {
		return nil
	}

	if sort.Field == "" {
		return errors.New("sort field must not be empty")
	}
	if sort.Order == "" {
		return errors.New("sort order must not be empty")
	}

	validFields := []any{SortFieldDefault, SortFieldID, SortFieldCreatedAt, SortFieldUpdatedAt, SortFieldTags, SortFieldLastRevisedAt, SortFieldNextRevisionAt}

	return validation.ValidateStruct(sort,
		validation.Field(&sort.Field, validation.Required, validation.In(validFields...).Error(fmt.Sprintf("sort field must be one of %v", validFields))),
		validation.Field(&sort.Order, validation.Required, validation.In(SortOrderDefault, SortOrderAsc, SortOrderDesc).Error("sort order must be one of asc, desc")),
	)

}
