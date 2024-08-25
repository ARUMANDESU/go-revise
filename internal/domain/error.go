package domain

type Error struct {
	Op      string
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Op + ": " + e.Message + ": " + e.Err.Error()
}

// Implement the Unwrap method to support errors.Is and errors.As
func (e *Error) Unwrap() error {
	return e.Err
}

func WrapError(err error, message string) error {
	return &Error{
		Message: message,
		Err:     err,
	}
}

func WrapErrorWithOp(err error, op, message string) error {
	return &Error{
		Op:      op,
		Message: message,
		Err:     err,
	}
}
