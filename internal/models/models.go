package models

import "time"

type User struct {
	UserID      string    `json:"user_id" db:"user_id"`
	Username    string    `json:"username" db:"username"`
	DisplayName string    `json:"display_name" db:"display_name"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	TeamName    *string   `json:"team_name" db:"team_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Team struct {
	TeamName  string    `json:"team_name" db:"team_name"`
	Desc      *string   `json:"description" db:"description"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
