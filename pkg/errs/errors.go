package errs

import (
	"errors"
	"fmt"
)

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
	ErrorTypeNotFound       = ErrorType{"not-found"}
	ErrorTypeConflict       = ErrorType{"conflict"}
)

type MsgError struct {
	err     error
	op      string
	msg     string
	errType ErrorType
}

func (s MsgError) Error() string {
	if s.err != nil {
		return fmt.Sprintf("op: %s, type: %s, msg: %s, error: %v", s.op, s.errType.t, s.msg, s.err)
	}
	return fmt.Sprintf("op: %s, type: %s", s.op, s.errType.t)
}

func (s MsgError) Message() string {
	return fmt.Sprintf("%s: %s: %v", s.errType.t, s.msg, s.err)
}

// Unwrap provides compatibility for Go 1.13+ error wrapping.
// It returns the underlying error, allowing for inspection of wrapped errors.
func (e *MsgError) Unwrap() error {
	return e.err
}

// Is checks if the error matches a target error, particularly useful for comparing types.
func (e *MsgError) Is(target error) bool {
	t, ok := target.(*MsgError)
	if !ok {
		return false
	}
	return e.errType.t == t.errType.t
}

// As allows the error to be cast to a target type.
func (e *MsgError) As(target interface{}) bool {
	if t, ok := target.(**MsgError); ok {
		*t = e
		return true
	}
	return errors.As(e.err, target)
}

func NewMsgError(op string, err error, msg string) MsgError {
	return MsgError{
		op:      op,
		err:     err,
		msg:     msg,
		errType: ErrorTypeUnknown,
	}
}

func NewAuthorizationError(op string, err error, msg string) MsgError {
	return MsgError{
		op:      op,
		err:     err,
		msg:     msg,
		errType: ErrorTypeAuthorization,
	}
}

func NewIncorrectInputError(op string, err error, msg string) MsgError {
	return MsgError{
		op:      op,
		err:     err,
		msg:     msg,
		errType: ErrorTypeIncorrectInput,
	}
}

func NewNotFoundError(op string, err error, msg string) MsgError {
	return MsgError{
		op:      op,
		err:     err,
		msg:     msg,
		errType: ErrorTypeNotFound,
	}
}

func NewConflictError(op string, err error, msg string) MsgError {
	return MsgError{
		op:      op,
		err:     err,
		msg:     msg,
		errType: ErrorTypeConflict,
	}
}
