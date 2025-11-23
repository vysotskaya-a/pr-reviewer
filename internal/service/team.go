package service

import (
	"context"
	"errors"
	"pr-reviewer/internal/models"
	"pr-reviewer/internal/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, description *string) (*models.Team, error)
	GetTeam(ctx context.Context, name string) (*models.Team, error)
	ListTeams(ctx context.Context) ([]models.Team, error)
	DeleteTeam(ctx context.Context, name string) error
}

var ErrTeamHasMembers = errors.New("team has members")

type teamService struct {
	teams repository.TeamRepository
}

func NewTeamService(t repository.TeamRepository) TeamService {
	return &teamService{teams: t}
}

func (s *teamService) CreateTeam(ctx context.Context, teamName string, description *string) (*models.Team, error) {
	return s.teams.Create(ctx, teamName, description)
}

func (s *teamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	return s.teams.GetByName(ctx, name)
}

func (s *teamService) ListTeams(ctx context.Context) ([]models.Team, error) {
	return s.teams.List(ctx)
}

func (s *teamService) DeleteTeam(ctx context.Context, name string) error {
	err := s.teams.Delete(ctx, name)
	if err != nil {
		if errors.Is(err, repository.ErrForeignKeyViolation) {
			return ErrTeamHasMembers
		}
		return err
	}
	return nil
}
