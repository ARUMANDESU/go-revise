package service

import "errors"

var (
	ErrInternal = errors.New("internal error")

	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrAlreadyExists   = errors.New("already exists")
)
