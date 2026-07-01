package apperr

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrAliasAlreadyExists = errors.New("alias already exists")
)
