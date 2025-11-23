package repository

import (
	"context"
	"time"

	"pr-reviewer/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepoPG struct {
	p *pgxpool.Pool
}

func NewUserRepositoryPG(p *pgxpool.Pool) UserRepository {
	return &userRepoPG{p: p}
}

func (r *userRepoPG) Create(ctx context.Context, username string, displayName *string, teamName *string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	query := `INSERT INTO users(username, display_name, team_name) VALUES ($1,$2,$3)
	          RETURNING user_id, username, display_name, is_active, team_name, created_at`
	var u models.User
	row := r.p.QueryRow(ctx, query, username, displayName, teamName)
	if err := row.Scan(&u.UserID, &u.Username, &u.DisplayName, &u.IsActive, &u.TeamName, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepoPG) GetByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var u models.User
	row := r.p.QueryRow(ctx, `SELECT user_id, username, display_name, is_active, team_name, created_at FROM users WHERE user_id = $1`, id)
	if err := row.Scan(&u.UserID, &u.Username, &u.DisplayName, &u.IsActive, &u.TeamName, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepoPG) List(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.p.Query(ctx, `SELECT user_id, username, display_name, is_active, team_name, created_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.DisplayName, &u.IsActive, &u.TeamName, &u.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func (r *userRepoPG) Update(ctx context.Context, id string, displayName *string, isActive *bool, teamName *string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.p.Exec(ctx, `UPDATE users SET display_name = COALESCE($1, display_name), is_active = COALESCE($2, is_active), team_name = COALESCE($3, team_name) WHERE user_id = $4`, displayName, isActive, teamName, id)
	return err
}

func (r *userRepoPG) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.p.Exec(ctx, `DELETE FROM users WHERE user_id = $1`, id)
	return err
}
