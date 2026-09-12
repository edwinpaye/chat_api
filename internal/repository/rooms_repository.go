package repository

import (
	"context"
	"database/sql"
	"errors"

	"chat_api/internal/domain"
)

// RoomsRepository is the persistence port for Rooms.
type RoomsRepository interface {
	Create(ctx context.Context, e *domain.Rooms) error
	Get(ctx context.Context, id int64) (*domain.Rooms, error)
	List(ctx context.Context, limit, offset int) ([]domain.Rooms, error)
	Update(ctx context.Context, e *domain.Rooms) error
	Delete(ctx context.Context, id int64) error
}

type roomsRepo struct{ db *sql.DB }

func NewRoomsRepository(db *sql.DB) RoomsRepository {
	return &roomsRepo{db: db}
}

const roomsColumns = "id, name"

func (r *roomsRepo) Create(ctx context.Context, e *domain.Rooms) error {
	q := `INSERT INTO rooms (name)
	VALUES (?)`
	res, err := r.db.ExecContext(ctx, q, e.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		e.Id = int64(id)
	}
	return nil}

func (r *roomsRepo) Get(ctx context.Context, id int64) (*domain.Rooms, error) {
	q := `SELECT ` + roomsColumns + ` FROM rooms WHERE id = ?`
	var e domain.Rooms
	err := r.db.QueryRowContext(ctx, q, id).Scan(&e.Id, &e.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *roomsRepo) List(ctx context.Context, limit, offset int) ([]domain.Rooms, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	q := `SELECT ` + roomsColumns + ` FROM rooms ORDER BY id LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Rooms
	for rows.Next() {
		var e domain.Rooms
		if err := rows.Scan(&e.Id, &e.Name); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *roomsRepo) Update(ctx context.Context, e *domain.Rooms) error {
	q := `UPDATE rooms SET name = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, e.Name, e.Id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *roomsRepo) Delete(ctx context.Context, id int64) error {
	q := `DELETE FROM rooms WHERE id = ?`
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