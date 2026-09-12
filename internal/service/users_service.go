package service

import (
	"context"
	"encoding/json"
	"errors"

	"chat_api/internal/domain"
	"chat_api/internal/repository"
)

// UsersService implements use cases for Users.
type UsersService struct {
	repo repository.UsersRepository
}

func newUsersService(r repository.UsersRepository) *UsersService {
	return &UsersService{repo: r}
}

func (s *UsersService) Create(ctx context.Context, payload json.RawMessage) (*domain.Users, error) {
	var e domain.Users
	if err := json.Unmarshal(payload, &e); err != nil {
		return nil, errors.Join(domain.ErrValidation, err)
	}
	if err := s.repo.Create(ctx, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func (s *UsersService) Get(ctx context.Context, id string) (*domain.Users, error) {
	return s.repo.Get(ctx, id)
}

func (s *UsersService) List(ctx context.Context, limit, offset int) ([]domain.Users, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *UsersService) Update(ctx context.Context, id string, payload json.RawMessage) (*domain.Users, error) {
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

func (s *UsersService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}