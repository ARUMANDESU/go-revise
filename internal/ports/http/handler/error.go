package handler

import (
	"errors"
	"net/http"

	"github.com/ARUMANDESU/go-revise/pkg/contexts"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	op := errs.Op("handler.handleError")
	log := contexts.Logger(r.Context())
	requestID := contexts.RequestID(r.Context())

	log.With("request_id", requestID)

	var appErr *errs.Error
	if !errors.As(err, &appErr) {
		appErr = errs.NewUnknownError(op, err, "unknown error")
	}
	appErr.Log(log)

	var statusCode int
	switch {
	case errs.IsErrorType(appErr, errs.ErrorTypeIncorrectInput):
		statusCode = http.StatusBadRequest
	case errs.IsErrorType(appErr, errs.ErrorTypeNotFound):
		statusCode = http.StatusNotFound
	case errs.IsErrorType(appErr, errs.ErrorTypeConflict),
		errs.IsErrorType(appErr, errs.ErrorTypeAlreadyExists):
		statusCode = http.StatusConflict
	case errs.IsErrorType(appErr, errs.ErrorTypeAuthorization):
		statusCode = http.StatusUnauthorized
	case errs.IsErrorType(appErr, errs.ErrorTypeForbidden):
		statusCode = http.StatusForbidden
	case errs.IsErrorType(appErr, errs.ErrorTypeUnknown):
		statusCode = http.StatusInternalServerError
	default:
		statusCode = http.StatusInternalServerError
	}

	response := map[string]any{
		"error":      appErr.Type().String(),
		"message":    appErr.Message()["message"],
		"request_id": requestID,
		"succeeded":  false,
	}

	err = writeJSON(w, statusCode, response, nil)
	if err != nil {
		log.Error("failed to write response", logutil.Err(err))
		return
	}
}
