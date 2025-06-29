package validation

import (
	"net/mail"
	"strings"
	"time"
	"users-backend/pkg/models"
)

func ValidateUser(name string, email string, birthday time.Time) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if err := ValidateEmail(email); err != nil {
		return err
	}
	if err := ValidateDateOfBirth(birthday); err != nil {
		return err
	}
	return nil
}

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return models.NewInternalError(
			models.ContextBadRequest,
			"email is required and cannot be empty",
		)
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return models.NewInternalError(
			models.ContextBadRequest,
			"email format is invalid",
		)
	}
	return nil
}

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.NewInternalError(
			models.ContextBadRequest,
			"name is required and cannot be empty",
		)
	}
	if len(name) < 2 || len(name) > 100 {
		return models.NewInternalError(
			models.ContextBadRequest,
			"name must be between 2 and 100 characters",
		)
	}
	return nil
}

func ValidateDateOfBirth(birthday time.Time) error {
	now := time.Now()
	if birthday.After(now) {
		return models.NewInternalError(
			models.ContextBadRequest,
			"date of birth cannot be in the future",
		)
	}
	return nil
}
