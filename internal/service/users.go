package service

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"Starostina-elena/pull_req_assign/internal/storage"
	"context"
	"errors"
	"log/slog"
)

type UserService interface {
	CreateUser(ctx context.Context, id string, name string, isActive bool) error
	GetUser(ctx context.Context, id string) (core.User, error)
	SetTeamToUser(ctx context.Context, userId string, teamId int64) error
	ExpellUserFromTeam(ctx context.Context, userId string) error
	SetUserIsActive(ctx context.Context, userId string, isActive bool) error
	GetPullRequestAssigned(ctx context.Context, userId string) ([]core.PullRequest, error)
}

type UserServiceImpl struct {
	storage *storage.DB
	log     *slog.Logger
}

func NewUserService(log *slog.Logger, storage *storage.DB) UserService {
	return &UserServiceImpl{storage: storage, log: log}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, id string, name string, isActive bool) error {
	if name == "" {
		return errors.New("name required")
	}
	return s.storage.AddUser(ctx, id, name, isActive)
}

func (s *UserServiceImpl) GetUser(ctx context.Context, id string) (core.User, error) {
	user, err := s.storage.GetUserByID(ctx, id)
	if err != nil {
		return core.User{}, err
	}
	return user, nil
}

func (s *UserServiceImpl) SetTeamToUser(ctx context.Context, userId string, teamId int64) error {
	return s.storage.SetTeamToUser(ctx, userId, teamId)
}

func (s *UserServiceImpl) ExpellUserFromTeam(ctx context.Context, userId string) error {
	return s.storage.ExpellUserFromTeam(ctx, userId)
}

func (s *UserServiceImpl) SetUserIsActive(ctx context.Context, userId string, isActive bool) error {
	return s.storage.SetUserIsActive(ctx, userId, isActive)
}

func (s *UserServiceImpl) GetPullRequestAssigned(ctx context.Context, userId string) ([]core.PullRequest, error) {
	return s.storage.GetPullRequestsAssigned(ctx, userId)
}
