package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"users-backend/pkg/models"
	"users-backend/pkg/services"

	"github.com/google/uuid"
)

type UserRequest struct {
	ID          uuid.UUID `json:"external_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"external_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

func NewUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		DateOfBirth: user.DateOfBirth,
	}
}

func (s *APIServer) HandleCreateUser(w http.ResponseWriter, r *http.Request) error {
	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return models.NewWrappedError(err, models.ContextBadRequest, "request body contains malformed data")
	}

	serviceReq := services.UserCreationRequest{
		ID:          userRequest.ID,
		Name:        userRequest.Name,
		Email:       userRequest.Email,
		DateOfBirth: userRequest.DateOfBirth,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	user, err := s.service.CreateUser(ctx, serviceReq)
	if err != nil {
		return err
	}

	response := NewUserResponse(user)
	return ConstructSuccessResponse(w, http.StatusCreated, response)
}

func (s *APIServer) HandleGetUser(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return models.NewWrappedError(err, models.ContextBadRequest, fmt.Sprintf("UUID '%s' is not formatted correctly.", id))
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	user, err := s.service.GetUser(ctx, userUUID)
	if err != nil {
		return err
	}

	response := NewUserResponse(user)
	return ConstructSuccessResponse(w, http.StatusOK, response)
}
