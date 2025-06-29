package models

import (
	"errors"
	"fmt"
)

type WrappedError struct {
	Base    error
	Context string
	Reason  string
}

type InternalError struct {
	Context string
	Reason  string
}

var (
	ErrMethodNotSupported      = errors.New("HTTP method not supported")
	ErrDatabaseEnvConfigNotSet = errors.New("database environment configuration is not set")
)

var (
	ContextNotFound           = "resource_not_found"
	ContextMethodNotSupported = "method_not_supported"
	ContextInternalServer     = "internal_server_error"
	ContextConflictValue      = "value_conflict"
	ContextBadRequest         = "bad_request"
)

func (e *WrappedError) Error() string {
	return fmt.Sprintf("%s: %s", e.Context, e.Reason)
}

func NewWrappedError(base error, context string, reason string) *WrappedError {
	return &WrappedError{Base: base, Context: context, Reason: reason}
}

func (e *WrappedError) Unwrap() error {
	return e.Base
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("%s: %s", e.Context, e.Reason)
}

func NewInternalError(context string, reason string) *InternalError {
	return &InternalError{Context: context, Reason: reason}
}
