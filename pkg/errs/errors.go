package errs

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

var (
	ErrInvalidInput = errors.New("invalid input")
)

type Message struct {
	Key   string
	Value string
}

type ErrorType struct {
	t string
}

func (e ErrorType) String() string {
	return e.t
}

func (e ErrorType) Is(target ErrorType) bool {
	return e.t == target.t
}

// Error types not exposed to client
var (
	ErrorTypeUnknown = ErrorType{"unknown"}
)

// is Error type exposed to client
var (
	ErrorTypeForbidden      = ErrorType{"forbidden"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
	ErrorTypeNotFound       = ErrorType{"not-found"}
	ErrorTypeConflict       = ErrorType{"conflict"}
	ErrorTypeAlreadyExists  = ErrorType{"already-exists"}
)

// Op describes an operation, usually as the package and method,
// such as db.get_user.
type Op string

type OpMessage struct {
	Op      Op
	Message string
}

func (o OpMessage) String() string {
	return fmt.Sprintf("%s: %s", o.Op, o.Message)
}

type Error struct {
	err        error // The underlying error, if there is one
	errType    ErrorType
	opMessages []OpMessage
	// context mainly used for logging, do not put sensitive information here, but it's okay to put some context
	// such as request ID, user ID, etc.
	context map[string]any
	// messages is exposed to client, do not put sensitive information here
	messages map[string]string
}

func (e *Error) Error() string {
	strBuilder := strings.Builder{}
	strBuilder.WriteString("type: ")
	strBuilder.WriteString(e.errType.t)
	for _, opMsg := range e.opMessages {
		strBuilder.WriteString(opMsg.String())
		strBuilder.WriteString(", ")
	}
	for k, v := range e.context {
		strBuilder.WriteString(k)
		strBuilder.WriteString(": ")
		strBuilder.WriteString(fmt.Sprintf("%v", v))
		strBuilder.WriteString(", ")
	}
	strBuilder.WriteString(fmt.Sprintf("original error: %v", e.err))

	return strBuilder.String()
}

func (e *Error) Message() map[string]string {
	return e.messages
}

func (e *Error) Type() ErrorType {
	return e.errType
}

func (e *Error) Operations() []OpMessage {
	return e.opMessages
}

func (e *Error) HasOperation(op Op) bool {
	for _, opMsg := range e.opMessages {
		if opMsg.Op == op {
			return true
		}
	}
	return false
}

func (e *Error) Log(logger *slog.Logger) {
	logger = logger.With("op", e.opMessages, "type", e.errType)
	for k, v := range e.context {
		logger = logger.With(k, v)
	}

	logger.Error("error occurred", logutil.Err(e))
}

func (e *Error) Unwrap() error {
	return e.err
}
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return errors.Is(e.err, target)
	}
	return e.errType.Is(t.errType)
}

func (e *Error) WithContext(key string, value any) *Error {
	if e.context == nil {
		e.context = make(map[string]any)
	}
	e.context[key] = value
	return e
}

func (e *Error) WithMessages(messages []Message) *Error {
	if e.messages == nil {
		e.messages = make(map[string]string)
	}
	for _, msg := range messages {
		e.messages[msg.Key] = msg.Value
	}
	return e
}

func IsErrorType(err error, errType ErrorType) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	return e.errType == errType
}

func WithOp(op Op, err error, msg string) error {
	if err == nil {
		return nil
	}

	var e *Error
	if !errors.As(err, &e) {
		e = E(op, err, ErrorTypeUnknown, msg)
	} else {
		e.opMessages = append(e.opMessages, OpMessage{Op: op, Message: msg})
	}

	return e
}

func E(op Op, err error, errType ErrorType, message string) *Error {
	return &Error{
		err:        err,
		errType:    errType,
		opMessages: []OpMessage{{op, message}},
		messages:   make(map[string]string),
		context:    make(map[string]any),
	}
}

func NewUnknownError(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeUnknown, message)
}

func NewIncorrectInputError(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeIncorrectInput, message)
}

func NewConflictError(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeConflict, message)
}

func NewNotFound(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeNotFound, message)
}

func NewAlreadyExistsError(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeAlreadyExists, message)
}

func NewAuthorizationError(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeAuthorization, message)
}

func NewForbiddenError(op Op, err error, message string) *Error {
	return E(op, err, ErrorTypeForbidden, message)
}
