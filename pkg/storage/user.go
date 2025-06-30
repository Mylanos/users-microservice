package storage

import (
	"time"
	"users-backend/pkg/models"

	"github.com/google/uuid"
)

// UserDTO represents the database structure for users
type UserEntity struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Email       string    `gorm:"uniqueIndex;not null"`
	DateOfBirth time.Time `gorm:"type:date"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (UserEntity) TableName() string {
	return "users"
}

func (dto *UserEntity) ToModel() *models.User {
	return &models.User{
		ID:          dto.ID,
		Name:        dto.Name,
		Email:       dto.Email,
		DateOfBirth: dto.DateOfBirth,
	}
}

func (dto *UserEntity) FromModel(user *models.User) {
	dto.ID = user.ID
	dto.Name = user.Name
	dto.Email = user.Email
	dto.DateOfBirth = user.DateOfBirth
}
