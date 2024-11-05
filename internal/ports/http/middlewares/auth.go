package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	initdata "github.com/telegram-mini-apps/init-data-golang"

	"github.com/ARUMANDESU/go-revise/pkg/contexts"
	"github.com/ARUMANDESU/go-revise/pkg/env"
	"github.com/ARUMANDESU/go-revise/pkg/logutil"
)

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("authorization"))
		if header == "" {
			slog.Info("header is empty")
			return
		}
		authParts := strings.Split(header, " ")
		if len(authParts) != 2 {
			slog.Info("authParts are not two", slog.Any("authParts", authParts))
			// TODO: return authentication error
			return
		}
		authType := authParts[0]
		authData := authParts[1]

		switch authType {
		case "tma":
			expIn := time.Hour
			if m.EnvMode != env.Local {
				if err := initdata.Validate(authData, m.tmaAuthToken, expIn); err != nil {
					slog.Info("failed to validate tma auth data", logutil.Err(err))
					// TODO: return authentication error
					return
				}
			}

			initData, err := initdata.Parse(authData)
			if err != nil {
				slog.Info("failed to parse tma auth data", logutil.Err(err))
				// TODO: return authentication error
				return
			}

			slog.Info("got tma authorization header", slog.Any("initData", initData))

			r = r.WithContext(contexts.WithTMAInitData(r.Context(), initData))
		default:
			slog.Info("is not any of the authorization type", slog.String("type", authType))
		}
		next.ServeHTTP(w, r)
	})
}
