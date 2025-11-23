package service

import (
	"context"
	"errors"
	"pr-reviewer/internal/models"
	"pr-reviewer/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrNoAvailableReviewers = errors.New("no available reviewers")
	ErrReviewerNotInPR      = errors.New("reviewer not assigned")
	ErrCannotModifyMerged   = errors.New("cannot modify merged PR")
)

type PRService interface {
	CreatePR(ctx context.Context, name string, authorID string) (*models.PullRequest, []models.User, error)
	ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (*models.User, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequest, error)
	GetPR(ctx context.Context, id string) (*models.PullRequest, []models.User, error)
	ListByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequest, error)
}

type prService struct {
	prRepo   repository.PRRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPRService(pr repository.PRRepository, users repository.UserRepository, teams repository.TeamRepository) PRService {
	return &prService{prRepo: pr, userRepo: users, teamRepo: teams}
}

func (s *prService) CreatePR(ctx context.Context, name string, authorID string) (*models.PullRequest, []models.User, error) {
	author, err := s.userRepo.GetByID(ctx, authorID)
	if err != nil {
		return nil, nil, err
	}
	if author.TeamName == nil {
		return nil, nil, errors.New("author has no team")
	}

	teamUsers, err := s.userRepo.ListUsersByTeam(ctx, *author.TeamName)
	if err != nil {
		return nil, nil, err
	}

	// choose reviewers
	var reviewers []models.User
	for _, u := range teamUsers {
		if u.UserID == authorID {
			continue
		}
		if u.IsActive {
			reviewers = append(reviewers, u)
		}
		if len(reviewers) >= 2 {
			break
		}
	}

	// create PR
	pr := &models.PullRequest{
		PullRequestID:   uuid.New().String(),
		PullRequestName: name,
		AuthorID:        authorID,
		TeamName:        *author.TeamName,
		Status:          "OPEN",
	}

	err = s.prRepo.Create(ctx, pr)
	if err != nil {
		return nil, nil, err
	}

	// assign reviewers
	for _, r := range reviewers {
		_ = s.prRepo.AddReviewer(ctx, pr.PullRequestID, r.UserID)
	}

	return pr, reviewers, nil
}

func (s *prService) GetPR(ctx context.Context, id string) (*models.PullRequest, []models.User, error) {
	pr, err := s.prRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	revs, err := s.prRepo.ListReviewers(ctx, id)
	return pr, revs, err
}

func (s *prService) ListByReviewer(ctx context.Context, reviewerID string) ([]models.PullRequest, error) {
	return s.prRepo.ListByReviewer(ctx, reviewerID)
}

func (s *prService) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return pr, nil
	}
	err = s.prRepo.SetMerged(ctx, prID)
	if err != nil {
		return nil, err
	}
	return s.prRepo.GetByID(ctx, prID)
}

func (s *prService) ReassignReviewer(ctx context.Context, prID string, oldReviewerID string) (*models.User, error) {
	pr, err := s.prRepo.GetByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return nil, ErrCannotModifyMerged
	}

	reviewers, err := s.prRepo.ListReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	// check old reviewer exists
	found := false
	for _, r := range reviewers {
		if r.UserID == oldReviewerID {
			found = true
		}
	}
	if !found {
		return nil, ErrReviewerNotInPR
	}

	// find new reviewer
	teamUsers, err := s.userRepo.ListUsersByTeam(ctx, pr.TeamName)
	if err != nil {
		return nil, err
	}

	var candidates []models.User
	for _, u := range teamUsers {
		if u.UserID == oldReviewerID {
			continue
		}
		if !u.IsActive {
			continue
		}
		alreadyAssigned := false
		for _, r := range reviewers {
			if r.UserID == u.UserID {
				alreadyAssigned = true
				break
			}
		}
		if !alreadyAssigned {
			candidates = append(candidates, u)
		}
	}

	if len(candidates) == 0 {
		return nil, ErrNoAvailableReviewers
	}

	newReviewer := candidates[0]

	_ = s.prRepo.RemoveReviewer(ctx, prID, oldReviewerID)
	_ = s.prRepo.AddReviewer(ctx, prID, newReviewer.UserID)

	return &newReviewer, nil
}
