package reviseitem

import (
	"time"

	"github.com/ARUMANDESU/go-revise/internal/domain/valueobject"
)

func validateName(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	return nil
}

func validateDescription(description string) error {
	if description == "" {
		return ErrDescriptionRequired
	}
	return nil
}

func validateTags(tags valueobject.StringArray) error {
	if !tags.IsValid() {
		return ErrInvalidTags
	}
	return nil
}

func validateNextRevisionAt(nextRevisionAt time.Time) error {
	if nextRevisionAt.IsZero() || nextRevisionAt.Before(time.Now()) {
		return ErrInvalidNextRevisionAt
	}
	return nil
}
