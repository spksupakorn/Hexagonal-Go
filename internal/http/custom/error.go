package custom

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string, defaultMsg string) error {
	if message == "" {
		message = defaultMsg
	}
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewBadRequestError(message string) error {
	return NewAppError(http.StatusBadRequest, message, "bad request")
}

func NewNotFoundError(message string) error {
	return NewAppError(http.StatusNotFound, message, "not found")
}

func NewUnauthorizedError(message string) error {
	return NewAppError(http.StatusUnauthorized, message, "unauthorized")
}

func NewValidationError(message string) error {
	return NewAppError(http.StatusUnprocessableEntity, message, "validation error")
}

func NewUnexpectedError(message string) error {
	return NewAppError(http.StatusInternalServerError, message, "unexpected error")
}

func NewForbiddenError(message string) error {
	return NewAppError(http.StatusForbidden, message, "forbidden")
}

func NewConflictError(message string) error {
	return NewAppError(http.StatusConflict, message, "conflict")
}

func NewNoContentError() error {
	return NewAppError(http.StatusNoContent, "no content", "no content")
}
