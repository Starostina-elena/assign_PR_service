package service

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"Starostina-elena/pull_req_assign/internal/storage"
	"context"
	"errors"
	"log/slog"
)

type UserService interface {
	CreateUser(ctx context.Context, name string, isActive bool) (int64, error)
	GetUser(ctx context.Context, id int64) (core.User, error)
}

type UserServiceImpl struct {
	storage *storage.DB
	log     *slog.Logger
}

func NewUserService(log *slog.Logger, storage *storage.DB) UserService {
	return &UserServiceImpl{storage: storage, log: log}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, name string, isActive bool) (int64, error) {
	if name == "" {
		return 0, errors.New("name required")
	}
	return s.storage.AddUser(ctx, name, isActive)
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id int64) (core.User, error) {
	user, err := s.storage.GetUserByID(ctx, id)
	if err != nil {
		return core.User{}, err
	}
	return user, nil
}
