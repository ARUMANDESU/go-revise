package tgboterr

import (
	"errors"
	"log/slog"
	"strings"

	tb "gopkg.in/telebot.v4"

	"github.com/ARUMANDESU/go-revise/pkg/errs"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

func OnError(err error, c tb.Context) {
	op := errs.Op("tgboterr.on_error")
	log := slog.Default()

	var appErr *errs.Error
	if !errors.As(err, &appErr) {
		appErr = errs.NewUnknownError(op, err, "unknown error")
	}
	appErr.Log(log)

	switch {
	case errs.IsErrorType(appErr, errs.ErrorTypeIncorrectInput):
		sendError(c, "incorrect input", appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeNotFound):
		sendError(c, "not found", appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeConflict):
		sendError(c, "conflict", appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeAlreadyExists):
		sendError(c, "already exists", appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeAuthorization):
		sendError(c, "authorization", appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeForbidden):
		sendError(c, "forbidden", appErr.Message()["message"])
	case errs.IsErrorType(appErr, errs.ErrorTypeUnknown):
		sendError(c, "internal error", "")
	default:
		sendError(c, "internal error", "")
	}
}

func sendError(c tb.Context, errType, msg string) {
	errorMessage := strings.Builder{}
	errorMessage.WriteString("😕 Oops! ")
	if msg != "" {
		errorMessage.WriteString(msg)
	} else {
		errorMessage.WriteString(getDefaultMessage(errType))
	}
	errorMessage.WriteString("\n\nTry: ")
	errorMessage.WriteString(getSuggestion(errType))
	err := c.Send(errorMessage.String())
	if err != nil {
		slog.Error(
			"failed to send error message",
			logutil.Err(err),
			slog.String("message", errorMessage.String()),
		)
	}
}

func getDefaultMessage(errType string) string {
	switch errType {
	case "incorrect input":
		return "The provided input format is incorrect"
	case "not found":
		return "We couldn't find what you're looking for"
	case "conflict":
		return "There's a conflict with existing data"
	case "already exists":
		return "This item already exists"
	case "authorization":
		return "You need to be authorized for this action"
	case "forbidden":
		return "You don't have permission for this action"
	default:
		return "Something went wrong on our end"
	}
}

func getSuggestion(errType string) string {
	switch errType {
	case "incorrect input":
		return "Double-check your input format and try again"
	case "not found":
		return "Verify the ID or search term and try again"
	case "conflict":
		return "Review your data and try a different value"
	case "already exists":
		return "Use a different identifier or update the existing item"
	case "authorization":
		return "Log in and try again"
	case "forbidden":
		return "Contact an administrator for access"
	default:
		return "Please try again later or contact support"
	}
}
