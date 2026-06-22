package repository

import "errors"

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
	ErrNotFound           = errors.New("entity not found")
	ErrUnauthorized       = errors.New("unauthorized")
)
