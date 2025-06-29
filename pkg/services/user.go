package services

import (
	"context"
	"log"
	"time"
	"users-backend/pkg/models"
	"users-backend/pkg/storage"
	"users-backend/pkg/validation"

	"github.com/google/uuid"
)

type UserService interface {
	GetUser(context.Context, uuid.UUID) (*models.User, error)
	CreateUser(context.Context, UserCreationRequest) (*models.User, error)
}

type UserCreationRequest struct {
	ID          uuid.UUID
	Name        string
	Email       string
	DateOfBirth time.Time
}

type userService struct {
	storage storage.Storage
}

func NewUserService(storage storage.Storage) (UserService, error) {
	return &userService{storage: storage}, nil
}

func (us *userService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := us.storage.RetrieveUser(id)
	if err != nil {
		return nil, err
	}

	us.logUserAccess(id)
	return user, nil
}

func (us *userService) CreateUser(ctx context.Context, req UserCreationRequest) (*models.User, error) {
	if err := validation.ValidateUser(req.Name, req.Email, req.DateOfBirth); err != nil {
		return nil, err
	}

	//some business
	if !us.isEligibleForRegistration(req.DateOfBirth) {
		return nil, models.NewInternalError(models.ContextBadRequest, "user must have atleast 13 years to register")
	}

	newUser := models.NewUser(req.ID, req.Name, req.Email, req.DateOfBirth)

	//store
	if err := us.storage.CreateUser(newUser); err != nil {
		return nil, err
	}

	us.logUserCreated(newUser.ID)
	return newUser, nil
}

func (us *userService) logUserAccess(id uuid.UUID) {
	log.Printf("User %s accessed at %v", id, time.Now())
}

func (us *userService) logUserCreated(id uuid.UUID) {
	log.Printf("User %s created at %v", id, time.Now())
}

func (us *userService) isEligibleForRegistration(dateOfBirth time.Time) bool {
	age := time.Since(dateOfBirth).Hours() / 24 / 365.25
	return age >= 13
}
