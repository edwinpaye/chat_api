package repository

import (
	"context"
	"database/sql"
	"errors"

	"chat_api/internal/domain"
)

// MessagesRepository is the persistence port for Messages.
type MessagesRepository interface {
	Create(ctx context.Context, e *domain.Messages) error
	Get(ctx context.Context, id int64) (*domain.Messages, error)
	List(ctx context.Context, limit, offset int) ([]domain.Messages, error)
	Update(ctx context.Context, e *domain.Messages) error
	Delete(ctx context.Context, id int64) error
}

type messagesRepo struct{ db *sql.DB }

func NewMessagesRepository(db *sql.DB) MessagesRepository {
	return &messagesRepo{db: db}
}

const messagesColumns = "id, room_id, author_id, body, status, created_at, updated_at"

func (r *messagesRepo) Create(ctx context.Context, e *domain.Messages) error {
	q := `INSERT INTO messages (room_id, author_id, body, status, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, q, e.RoomId, e.AuthorId, e.Body, e.Status, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err == nil {
		e.Id = int64(id)
	}
	return nil}

func (r *messagesRepo) Get(ctx context.Context, id int64) (*domain.Messages, error) {
	q := `SELECT ` + messagesColumns + ` FROM messages WHERE id = ?`
	var e domain.Messages
	err := r.db.QueryRowContext(ctx, q, id).Scan(&e.Id, &e.RoomId, &e.AuthorId, &e.Body, &e.Status, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *messagesRepo) List(ctx context.Context, limit, offset int) ([]domain.Messages, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	q := `SELECT ` + messagesColumns + ` FROM messages ORDER BY id LIMIT ? OFFSET ?`
	rows, err := r.db.QueryContext(ctx, q, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Messages
	for rows.Next() {
		var e domain.Messages
		if err := rows.Scan(&e.Id, &e.RoomId, &e.AuthorId, &e.Body, &e.Status, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *messagesRepo) Update(ctx context.Context, e *domain.Messages) error {
	q := `UPDATE messages SET room_id = ?, author_id = ?, body = ?, status = ?, created_at = ?, updated_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, q, e.RoomId, e.AuthorId, e.Body, e.Status, e.CreatedAt, e.UpdatedAt, e.Id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *messagesRepo) Delete(ctx context.Context, id int64) error {
	q := `DELETE FROM messages WHERE id = ?`
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