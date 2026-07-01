package apperr

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}

// IsValidationError проверяет, является ли ошибка ошибкой валидации
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
