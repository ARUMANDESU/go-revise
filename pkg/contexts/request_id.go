package contexts

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
)

func GetRequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}
