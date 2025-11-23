package repository

import (
	"context"
	"pr-reviewer/internal/models"
)

type TeamRepository interface {
	Create(ctx context.Context, teamName string, description *string) (*models.Team, error)
	GetByName(ctx context.Context, name string) (*models.Team, error)
	List(ctx context.Context) ([]models.Team, error)
	Delete(ctx context.Context, name string) error
}
