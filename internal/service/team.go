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
	AttachUser(ctx context.Context, teamName string, userID *string, username string, isActive bool) error
}

var ErrTeamHasMembers = errors.New("team has members")

type teamService struct {
	teams repository.TeamRepository
	users repository.UserRepository
}

func NewTeamService(t repository.TeamRepository, u repository.UserRepository) TeamService {
	return &teamService{teams: t, users: u}
}

func (s *teamService) AttachUser(
	ctx context.Context,
	teamName string,
	userID *string,
	username string,
	isActive bool,
) error {

	// Проверяем что команда существует
	_, err := s.teams.GetByName(ctx, teamName)
	if err != nil {
		return errors.New("team not found")
	}

	// CASE 1: user_id не передан → создаём нового юзера
	if userID == nil {
		display := username
		_, err = s.users.Create(ctx, username, &display, &teamName)
		return err
	}

	// CASE 2: user_id передан — обновляем существующего
	_, err = s.users.GetByID(ctx, *userID)
	if err != nil {
		return err
	}

	display := username
	return s.users.Update(ctx, *userID, &display, &isActive, &teamName)
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
