//nolint:typecheck
package models

import (
	"errors"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/google/uuid"
)

type Event struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Date  time.Time `json:"date"`
	User  int       `json:"user"`
}

func (e *Event) Validate(requiredFilds []string) error {
	var (
		emptyUUID uuid.UUID
		errs      *multierror.Error
	)

	for _, fieldName := range requiredFilds {
		switch fieldName {
		case "id":
			if e.ID == emptyUUID {
				errs = multierror.Append(errs, errors.New("поле id не может быть пустым"))
			}
		case "title":
			if e.Title == "" {
				errs = multierror.Append(errs, errors.New("поле title не может быть пустым"))
			}
		case "date":
			if e.Date.IsZero() {
				errs = multierror.Append(errs, errors.New("поле date не может быть пустым"))
			}
		case "user":
			if e.User == 0 {
				errs = multierror.Append(errs, errors.New("поле user не может быть пустым или равно 0"))
			}
		}
	}

	return errs.ErrorOrNil() //nolint:wrapcheck
}
