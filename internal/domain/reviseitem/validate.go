package reviseitem

import (
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/multierr"
)

var (
	ErrInvalidID            = fmt.Errorf("invalid revise item id")
	ErrInvalidUserID        = fmt.Errorf("invalid userID")
	ErrInvalidArgument      = fmt.Errorf("invalid argument")
	ErrNextRevisionAtInPast = fmt.Errorf("nextRevisionAt must be in the future")
)

const (
	maxNameLength        = 255
	maxDescriptionLength = 1024
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

func validateNextRevisionAt(nextRevisionAt time.Time) error {
	if nextRevisionAt.IsZero() {
		return multierr.Combine(ErrInvalidArgument, fmt.Errorf("nextRevisionAt is required"))
	}
	if nextRevisionAt.Before(time.Now()) {
		return multierr.Combine(ErrInvalidArgument, ErrNextRevisionAtInPast)
	}

	return nil
}
