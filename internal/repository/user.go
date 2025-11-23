package repository

import (
	"context"
	"pr-reviewer/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, username string, displayName *string, teamName *string) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, id string, displayName *string, isActive *bool, teamName *string) error
	Delete(ctx context.Context, id string) error
}
