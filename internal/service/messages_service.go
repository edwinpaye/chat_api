package service

import (
	"context"
	"encoding/json"
	"errors"

	"chat_api/internal/domain"
	"chat_api/internal/repository"
)

// MessagesService implements use cases for Messages.
type MessagesService struct {
	repo repository.MessagesRepository
}

func newMessagesService(r repository.MessagesRepository) *MessagesService {
	return &MessagesService{repo: r}
}

func (s *MessagesService) Create(ctx context.Context, payload json.RawMessage) (*domain.Messages, error) {
	var e domain.Messages
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, errors.Join(domain.ErrValidation, err)
	}
	if err := s.repo.Create(ctx, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *MessagesService) Get(ctx context.Context, id int64) (*domain.Messages, error) {
	return s.repo.Get(ctx, id)
}

func (s *MessagesService) List(ctx context.Context, limit, offset int) ([]domain.Messages, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *MessagesService) Update(ctx context.Context, id int64, payload json.RawMessage) (*domain.Messages, error) {
	e, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(payload, e); err != nil {
		return nil, errors.Join(domain.ErrValidation, err)
	}
	e.Id = id
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *MessagesService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}