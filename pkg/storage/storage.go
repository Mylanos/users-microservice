package storage

import (
	"errors"
	"fmt"
	"strings"
	"users-microservice/pkg/config"
	"users-microservice/pkg/models"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	CreateUser(*models.User) error
	RetrieveUser(uuid.UUID) (*models.User, error)
	Close() error
}

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{TranslateError: false})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	db.AutoMigrate(&UserEntity{})

	return &PostgresStorage{db: db}, nil
}

func (ps *PostgresStorage) CreateUser(user *models.User) error {
	dto := &UserEntity{}
	dto.FromModel(user)

	tx := ps.db.Create(dto)
	if tx.Error != nil {
		errMsg := tx.Error.Error()
		// Check for PostgreSQL duplicate key constraint violations
		if strings.Contains(errMsg, "duplicate key value violates unique constraint") {
			if strings.Contains(errMsg, "users_pkey") {
				return models.NewWrappedError(tx.Error, models.ContextConflictValue, fmt.Sprintf("user with ID '%s' already exists", dto.ID))
			} else if strings.Contains(errMsg, "idx_users_email") {
				return models.NewWrappedError(tx.Error, models.ContextConflictValue, fmt.Sprintf("email '%s' is already in use", dto.Email))
			} else {
				// fallback
				return models.NewWrappedError(tx.Error, models.ContextConflictValue, "duplicate value violates unique constraint")
			}
		} else {
			return models.NewWrappedError(tx.Error, models.ContextInternalServer, "unexpected error while creating new user")
		}
	}
	return nil
}

func (ps *PostgresStorage) RetrieveUser(id uuid.UUID) (*models.User, error) {
	dto := &UserEntity{}
	tx := ps.db.First(dto, id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, models.NewWrappedError(tx.Error, models.ContextNotFound, fmt.Sprintf("user with '%s' ID does not exist", id))
		} else {
			return nil, models.NewWrappedError(tx.Error, models.ContextInternalServer, fmt.Sprintf("unexpected error while searching user with '%s' ID", id))
		}
	}
	return dto.ToModel(), nil
}

func (ps *PostgresStorage) CleanupTable() error {
	stmt := &gorm.Statement{DB: ps.db}
	if err := stmt.Parse(UserEntity{}); err != nil {
		return fmt.Errorf("failed to parse model for table name: %w", err)
	}
	tableName := stmt.Schema.Table

	switch ps.db.Dialector.Name() {
	case "postgres":
		if err := ps.db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", tableName)).Error; err != nil {
			return fmt.Errorf("failed to cleanup the table %s: %w", tableName, err)
		}
	}
	return nil
}

func (ps *PostgresStorage) Close() error {
	sqlDB, err := ps.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
