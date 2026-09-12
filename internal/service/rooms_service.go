package service

import (
	"context"
	"encoding/json"
	"errors"

	"chat_api/internal/domain"
	"chat_api/internal/repository"
)

// RoomsService implements use cases for Rooms.
type RoomsService struct {
	repo repository.RoomsRepository
}

func newRoomsService(r repository.RoomsRepository) *RoomsService {
	return &RoomsService{repo: r}
}

func (s *RoomsService) Create(ctx context.Context, payload json.RawMessage) (*domain.Rooms, error) {
	var e domain.Rooms
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, errors.Join(domain.ErrValidation, err)
	}
	if err := s.repo.Create(ctx, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *RoomsService) Get(ctx context.Context, id int64) (*domain.Rooms, error) {
	return s.repo.Get(ctx, id)
}

func (s *RoomsService) List(ctx context.Context, limit, offset int) ([]domain.Rooms, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *RoomsService) Update(ctx context.Context, id int64, payload json.RawMessage) (*domain.Rooms, error) {
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

func (s *RoomsService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}