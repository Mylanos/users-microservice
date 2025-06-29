package api

import (
	"errors"
	"net/http"
	"users-backend/pkg/models"
)

type APIError struct {
	Status      int    `json:"-"`
	Description string `json:"description"`
}

func NewAPIError(status int, description string) *APIError {
	return &APIError{Status: status, Description: description}
}

func (e *APIError) Error() string {
	return e.Description
}

func TranslateToAPIError(err error) APIError {
	var wrappedErr *models.WrappedError
	if errors.As(err, &wrappedErr) {
		return translateByContext(wrappedErr.Context, err.Error())
	}
	var internalErr *models.InternalError
	if errors.As(err, &internalErr) {
		return translateByContext(internalErr.Context, err.Error())
	}
	return *NewAPIError(http.StatusInternalServerError, err.Error())
}

func translateByContext(context string, message string) APIError {
	switch context {
	case models.ContextMethodNotSupported:
		return *NewAPIError(http.StatusMethodNotAllowed, message)
	case models.ContextConflictValue:
		return *NewAPIError(http.StatusConflict, message)
	case models.ContextNotFound:
		return *NewAPIError(http.StatusNotFound, message)
	case models.ContextBadRequest:
		return *NewAPIError(http.StatusBadRequest, message)
	default:
		return *NewAPIError(http.StatusInternalServerError, message)
	}
}
