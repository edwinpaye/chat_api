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

const usersColumns = "id, username, email, status, created_at, updated_at"

func (r *usersRepo) Create(ctx context.Context, e *domain.Users) error {
	q := `INSERT INTO users (username, email, status, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, q, e.Username, e.Email, e.Status, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		e.Id = string(id)
	}
	return nil}

func (r *usersRepo) Get(ctx context.Context, id string) (*domain.Users, error) {
	q := `SELECT ` + usersColumns + ` FROM users WHERE id = ?`
	var e domain.Users
	err := r.db.QueryRowContext(ctx, q, id).Scan(&e.Id, &e.Username, &e.Email, &e.Status, &e.CreatedAt, &e.UpdatedAt)
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
	q := `SELECT ` + usersColumns + ` FROM users ORDER BY id LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Users
	for rows.Next() {
		var e domain.Users
		if err := rows.Scan(&e.Id, &e.Username, &e.Email, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *usersRepo) Update(ctx context.Context, e *domain.Users) error {
	q := `UPDATE users SET username = ?, email = ?, status = ?, created_at = ?, updated_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, e.Username, e.Email, e.Status, e.CreatedAt, e.UpdatedAt, e.Id)
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
	q := `DELETE FROM users WHERE id = ?`
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