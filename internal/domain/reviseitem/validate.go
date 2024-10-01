package reviseitem

import (
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/multierr"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

var (
	ErrInvalidUserID        = fmt.Errorf("invalid userID")
	ErrInvalidArgument      = fmt.Errorf("invalid argument")
	ErrNextRevisionAtInPast = fmt.Errorf("nextRevisionAt must be in the future")
)

const (
	maxNameLength        = 255
	maxDescriptionLength = 1024
	maxNumTags           = 10
	maxTagLength         = 255
)

func validateName(name string) error {
	name = strings.TrimSpace(name)

	err := validation.Validate(name,
		validation.Required.Error("name is required"),
		validation.Length(1, maxNameLength).Error(fmt.Sprintf("name must be between 1 and %d characters", maxNameLength)),
	)
	if err != nil {
		return multierr.Combine(ErrInvalidArgument, err)
	}
	return nil
}

func validateDescription(description string) error {
	description = strings.TrimSpace(description)

	err := validation.Validate(description,
		validation.Length(1, maxDescriptionLength).Error(fmt.Sprintf("description must be between 1 and %d characters", maxDescriptionLength)),
	)
	if err != nil {
		return multierr.Combine(ErrInvalidArgument, err)
	}
	return nil
}

func validateTags(tags valueobject.StringArray) error {
	tags = tags.TrimSpace()
	err := validation.Validate(tags,
		validation.Length(0, maxNumTags).Error(fmt.Sprintf("tags must be between 0 and %d", maxNumTags)),
		validation.Each(validation.Length(1, maxTagLength).Error(fmt.Sprintf("tag must be between 1 and %d characters", maxTagLength))),
	)
	if err != nil {
		return multierr.Combine(ErrInvalidArgument, err)
	}
	return nil
}

func validateNextRevisionAt(nextRevisionAt time.Time) error {
	if nextRevisionAt.IsZero() {
		return multierr.Combine(ErrInvalidArgument, fmt.Errorf("nextRevisionAt is required"))
	}
	if nextRevisionAt.Before(time.Now()) {
		return multierr.Combine(ErrInvalidArgument, ErrNextRevisionAtInPast)
	}

	return nil
}
