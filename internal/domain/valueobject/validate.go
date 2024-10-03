package valueobject

import (
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	maxNumTags   = 10
	maxTagLength = 255
)

func ValidateTags(value any) error {
	tags, ok := value.(Tags)
	if !ok {
		return errors.New("invalid tags type")
	}

	err := validation.ValidateStruct(&tags,
		validation.Field(&tags.tags,
			validation.Length(0, maxNumTags).Error(fmt.Sprintf("max number of tags is %d", maxNumTags)),
			validation.Each(validation.Length(1, maxTagLength).Error(fmt.Sprintf("tag must be between 1 and %d characters", maxTagLength))),
		),
	)
	if err != nil {
		return err
	}

	return nil
}
