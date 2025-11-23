package models

import "time"

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PullRequest struct {
	PullRequestID   string     `json:"pull_request_id" db:"pull_request_id"`
	PullRequestName string     `json:"pull_request_name" db:"pull_request_name"`
	AuthorID        string     `json:"author_id" db:"author_id"`
	TeamName        string     `json:"team_name" db:"team_name"`
	Status          PRStatus   `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"createdat"`
	MergedAt        *time.Time `json:"merged_at,omitempty" db:"mergedat"`
}

type PRReviewer struct {
	PullRequestID string    `json:"pull_request_id" db:"pull_request_id"`
	ReviewerID    string    `json:"reviewer_id" db:"reviewer_id"`
	AssignedAt    time.Time `json:"assigned_at" db:"assigned_at"`
}
