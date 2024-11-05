package middlewares

import "github.com/ARUMANDESU/go-revise/pkg/env"

type Middleware struct {
	EnvMode env.Mode
	// telegram bot secret token
	tmaAuthToken string
}

func NewMiddleware(envMode env.Mode, botToken string) Middleware {
	return Middleware{
		EnvMode:      envMode,
		tmaAuthToken: botToken,
	}
}
