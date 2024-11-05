package contexts

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}
