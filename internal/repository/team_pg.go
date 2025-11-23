package repository

import (
	"context"
	"errors"
	"time"

	"pr-reviewer/internal/models"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type teamRepoPG struct {
	p *pgxpool.Pool
}

func NewTeamRepositoryPG(p *pgxpool.Pool) TeamRepository {
	return &teamRepoPG{p: p}
}

func (r *teamRepoPG) Create(ctx context.Context, teamName string, description *string) (*models.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var t models.Team
	err := r.p.QueryRow(ctx, `INSERT INTO teams(team_name, description) VALUES ($1,$2) RETURNING team_name, description, created_at`, teamName, description).
		Scan(&t.TeamName, &t.Desc, &t.CreatedAt)
	return &t, err
}

func (r *teamRepoPG) GetByName(ctx context.Context, name string) (*models.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	var t models.Team
	if err := r.p.QueryRow(ctx, `SELECT team_name, description, created_at FROM teams WHERE team_name = $1`, name).
		Scan(&t.TeamName, &t.Desc, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *teamRepoPG) List(ctx context.Context) ([]models.Team, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	rows, err := r.p.Query(ctx, `SELECT team_name, description, created_at FROM teams ORDER BY team_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.TeamName, &t.Desc, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

var ErrForeignKeyViolation = errors.New("foreign key violation")

func (r *teamRepoPG) Delete(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := r.p.Exec(ctx, `DELETE FROM teams WHERE team_name = $1`, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return ErrForeignKeyViolation
			}
		}
		return err
	}
	return nil
}
