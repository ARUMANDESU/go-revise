package httperr

import (
	"errors"
	"net/http"

	"github.com/ARUMANDESU/go-revise/internal/ports/http/httpio"
	"github.com/ARUMANDESU/go-revise/pkg/contexts"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	op := errs.Op("handler.handle_error")
	log := contexts.Logger(r.Context())
	requestID := contexts.RequestID(r.Context())

	log.With("request_id", requestID)

	var appErr *errs.Error
	if !errors.As(err, &appErr) {
		appErr = errs.NewUnknownError(op, err, "unknown error")
	}
	appErr.Log(log)

	switch {
	case errs.IsErrorType(appErr, errs.ErrorTypeIncorrectInput):
		BadRequest(w, r, appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeNotFound):
		NotFound(w, r, appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeConflict),
		errs.IsErrorType(appErr, errs.ErrorTypeAlreadyExists):
		Conflict(w, r, appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeAuthorization):
		Unauthorized(w, r, appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeForbidden):
		Forbidden(w, r, appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeUnknown):
		InternalServerError(w, r)
	default:
		InternalServerError(w, r)
	}
}

func Error(w http.ResponseWriter, r *http.Request, status int, errStr string, message string) {
	log := contexts.Logger(r.Context())
	requestID := contexts.RequestID(r.Context())

	log.With("request_id", requestID)
	response := map[string]any{
		"error":      errStr,
		"message":    message,
		"request_id": requestID,
		"succeeded":  false,
	}

	err := httpio.WriteJSON(w, status, response, nil)
	if err != nil {
		log.Error("failed to write response", logutil.Err(err))
	}
}

func Unauthorized(w http.ResponseWriter, r *http.Request, message string) {
	Error(w, r, http.StatusUnauthorized, "unauthorized", message)
}

func BadRequest(w http.ResponseWriter, r *http.Request, message string) {
	Error(w, r, http.StatusBadRequest, "bad-request", message)
}

func NotFound(w http.ResponseWriter, r *http.Request, message string) {
	Error(w, r, http.StatusNotFound, "not-found", message)
}

func Conflict(w http.ResponseWriter, r *http.Request, message string) {
	Error(w, r, http.StatusConflict, "conflict", message)
}

func Forbidden(w http.ResponseWriter, r *http.Request, message string) {
	Error(w, r, http.StatusForbidden, "forbidden", message)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	Error(w, r, http.StatusInternalServerError, "internal-server-error", "internal server error")
}
