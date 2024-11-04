package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ARUMANDESU/go-revise/internal/application"
	"github.com/ARUMANDESU/go-revise/internal/ports/http/handler"
	"github.com/ARUMANDESU/go-revise/pkg/errs"
)

type Port struct {
	server  *http.Server
	mux     *chi.Mux
	handler *handler.Handler
}

func NewHTTPPort(app application.Application) *Port {
	return &Port{
		handler: handler.NewHandler(app),
		mux:     chi.NewRouter(),
	}
}

// Start starts the http server.
//
//	NOTE: This function will block the current goroutine.
func (p *Port) Start(port string) error {
	op := errs.Op("http.Port.Start")
	p.setUpRouter()
	p.server = &http.Server{
		Addr:    ":" + port,
		Handler: p.mux,
	}

	slog.Info("http server started", slog.String("port", port))
	err := p.server.ListenAndServe()
	if err != nil {
		return errs.NewUnknownError(op, err, "failed to start http server")
	}
	return nil
}

func (p *Port) Stop() error {
	op := errs.Op("http.Port.Stop")
	if p.server == nil {
		return errs.NewUnknownError(op, nil, "server is not running")
	}

	err := p.server.Close()
	if err != nil {
		return errs.NewUnknownError(op, err, "failed to stop http server")
	}
	return nil
}
