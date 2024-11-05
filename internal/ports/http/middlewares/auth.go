package middlewares

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"

	"github.com/ARUMANDESU/go-revise/internal/ports/http/httperr"
	"github.com/ARUMANDESU/go-revise/pkg/contexts"
	"github.com/ARUMANDESU/go-revise/pkg/env"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

func (m *Middleware) Auth(next http.Handler) http.Handler {
	op := errs.Op("middleware.auth")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("authorization"))
		if header == "" {
			slog.Info("header is empty")
			return
		}
		authParts := strings.Split(header, " ")
		if len(authParts) != 2 {
			err := errs.
				NewIncorrectInputError(op, nil, "invalid authorization header format").
				WithMessages([]errs.Message{{Key: "message", Value: "invalid authorization header format, should be '<auth_type> <auth_data>'"}}).
				WithContext("header", header)
			httperr.HandleError(w, r, err)
			return
		}
		authType := authParts[0]
		authData := authParts[1]

		switch authType {
		case "tma":
			expIn := time.Hour
			if m.EnvMode != env.Local {
				if err := initdata.Validate(authData, m.tmaAuthToken, expIn); err != nil {
					handleInitDataError(w, r, err, "failed to validate tma authorization", authData)
					return
				}
			}

			initData, err := initdata.Parse(authData)
			if err != nil {
				handleInitDataError(w, r, err, "failed to parse tma authorization", authData)
				return
			}

			r = r.WithContext(contexts.WithTMAInitData(r.Context(), initData))
		default:
			err := errs.
				NewIncorrectInputError(op, nil, "unsupported authorization type").
				WithMessages([]errs.Message{{Key: "message", Value: "unsupported authorization type"}}).
				WithContext("auth_type", authType).
				WithContext("auth_data", authData)
			httperr.HandleError(w, r, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleInitDataError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	msg string,
	authData string,
) {
	op := errs.Op("middleware.auth.handle_init_data_error")
	var appErr *errs.Error
	switch {
	case errors.Is(err, initdata.ErrUnexpectedFormat):
		appErr = errs.
			NewAuthorizationError(op, err, "unexpected format").
			WithMessages([]errs.Message{{Key: "message", Value: "unexpected tma auth data format"}})
	case errors.Is(err, initdata.ErrAuthDateMissing):
		appErr = errs.
			NewAuthorizationError(op, err, "missing auth data").
			WithMessages([]errs.Message{{Key: "message", Value: "missing tma auth data"}})
	case errors.Is(err, initdata.ErrExpired):
		appErr = errs.
			NewAuthorizationError(op, err, "auth data expired").
			WithMessages([]errs.Message{{Key: "message", Value: "tma auth data expired"}})
	case errors.Is(err, initdata.ErrSignInvalid):
		appErr = errs.
			NewAuthorizationError(op, err, "invalid signature").
			WithMessages([]errs.Message{{Key: "message", Value: "tma auth data invalid signature"}})
	case errors.Is(err, initdata.ErrSignMissing):
		appErr = errs.
			NewAuthorizationError(op, err, "missing signature").
			WithMessages([]errs.Message{{Key: "message", Value: "tma auth data missing signature"}})
	default:
		appErr = errs.NewUnknownError(op, err, msg)
	}
	err = appErr.WithContext("auth_data", authData)

	httperr.HandleError(w, r, err)
}
