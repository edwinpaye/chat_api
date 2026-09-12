package repository

import (
	"context"
	"database/sql"
	"errors"

	"chat_api/internal/domain"
)

// UsersRepository is the persistence port for Users.
type UsersRepository interface {
	Create(ctx context.Context, e *domain.Users) error
	Get(ctx context.Context, id string) (*domain.Users, error)
	List(ctx context.Context, limit, offset int) ([]domain.Users, error)
	Update(ctx context.Context, e *domain.Users) error
	Delete(ctx context.Context, id string) error
}

type usersRepo struct{ db *sql.DB }

func NewUsersRepository(db *sql.DB) UsersRepository {
	return &usersRepo{db: db}
}

const usersColumns = "id, username, email, status, online, last_seen, created_at, updated_at"

func (r *usersRepo) Create(ctx context.Context, e *domain.Users) error {
	q := `INSERT INTO auth.users (username, email, status, online, last_seen, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
q += " RETURNING id"
	return r.db.QueryRowContext(ctx, q, e.Username, e.Email, e.Status, e.Online, e.LastSeen, e.CreatedAt, e.UpdatedAt).Scan(&e.Id)}

func (r *usersRepo) Get(ctx context.Context, id string) (*domain.Users, error) {
	q := `SELECT ` + usersColumns + ` FROM auth.users WHERE id = $1`
	var e domain.Users
	err := r.db.QueryRowContext(ctx, q, id).Scan(&e.Id, &e.Username, &e.Email, &e.Status, &e.Online, &e.LastSeen, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *usersRepo) List(ctx context.Context, limit, offset int) ([]domain.Users, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	q := `SELECT ` + usersColumns + ` FROM auth.users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Users
	for rows.Next() {
		var e domain.Users
		if err := rows.Scan(&e.Id, &e.Username, &e.Email, &e.Status, &e.Online, &e.LastSeen, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *usersRepo) Update(ctx context.Context, e *domain.Users) error {
	q := `UPDATE auth.users SET username = $1, email = $2, status = $3, online = $4, last_seen = $5, created_at = $6, updated_at = $7 WHERE id = $8`
	res, err := r.db.ExecContext(ctx, q, e.Username, e.Email, e.Status, e.Online, e.LastSeen, e.CreatedAt, e.UpdatedAt, e.Id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *usersRepo) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM auth.users WHERE id = $1`
	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}