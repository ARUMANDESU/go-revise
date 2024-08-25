package service

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	ValidTagsMinLength = 2
	ValidTagsMaxLength = 50

	ValidNameMinLength = 2
	ValidNameMaxLength = 50

	ValidDescriptionMinLength = 0
	ValidDescriptionMaxLength = 1000
)

func validateTags(value any) error {
	tags, ok := value.([]string)
	if !ok {
		return errors.New("tags must be a slice of strings")
	}

	return validation.Validate(tags, validation.Each(validation.Required, validation.Length(ValidTagsMinLength, ValidTagsMaxLength)))
}

func validateName(value any) error {
	name, ok := value.(string)
	if !ok {
		return errors.New("name must be a string")
	}

	return validation.Validate(name, validation.Length(ValidNameMinLength, ValidNameMaxLength))
}

func validateDescription(value any) error {
	description, ok := value.(string)
	if !ok {
		return errors.New("description must be a string")
	}

	return validation.Validate(description, validation.Length(ValidDescriptionMinLength, ValidDescriptionMaxLength))
}
