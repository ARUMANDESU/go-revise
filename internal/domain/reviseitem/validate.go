package reviseitem

import (
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.uber.org/multierr"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
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
	op := errs.Op("domain.reviseitem.validate_name")
	name = strings.TrimSpace(name)

	err := validation.Validate(name,
		validation.Required.Error("name is required"),
		validation.Length(1, maxNameLength).Error(fmt.Sprintf("name must be between 1 and %d characters", maxNameLength)),
	)
	if err != nil {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "invalid name").
			WithMessages([]errs.Message{{Key: "message", Value: err.Error()}}).
			WithContext("name", name)
	}
	return nil
}

func validateDescription(description string) error {
	op := errs.Op("domain.reviseitem.validate_description")
	description = strings.TrimSpace(description)

	err := validation.Validate(description,
		validation.Length(1, maxDescriptionLength).Error(fmt.Sprintf("description must be between 1 and %d characters", maxDescriptionLength)),
	)
	if err != nil {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "invalid description").
			WithMessages([]errs.Message{{Key: "message", Value: err.Error()}}).
			WithContext("description", description)
	}
	return nil
}

func validateNextRevisionAt(nextRevisionAt time.Time) error {
	op := errs.Op("domain.reviseitem.validate_next_revision_at")
	if nextRevisionAt.IsZero() {
		return errs.
			NewIncorrectInputError(op, errs.ErrInvalidInput, "nextRevisionAt is required").
			WithMessages([]errs.Message{{Key: "message", Value: "nextRevisionAt is required"}}).
			WithContext("nextRevisionAt", nextRevisionAt)
	}
	if nextRevisionAt.Before(time.Now()) {
		return multierr.Combine(ErrInvalidArgument, ErrNextRevisionAtInPast)
	}

	return nil
}
