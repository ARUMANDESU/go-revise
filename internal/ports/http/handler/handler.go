package handler

import (
	"github.com/ARUMANDESU/go-revise/internal/application"
)

type Handler struct {
	app application.Application
}

func NewHandler(app application.Application) *Handler {
	return &Handler{app: app}
}
