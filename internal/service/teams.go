package service

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"Starostina-elena/pull_req_assign/internal/storage"
	"context"
	"errors"
	"log/slog"
)

type TeamService interface {
	CreateTeam(ctx context.Context, name string) (int64, error)
	GetTeam(ctx context.Context, id int64) (core.Team, error)
	GetTeamMembers(ctx context.Context, teamId int64) ([]core.User, error)
}

type TeamServiceImpl struct {
	storage *storage.DB
	log     *slog.Logger
}

func NewTeamService(log *slog.Logger, storage *storage.DB) TeamService {
	return &TeamServiceImpl{storage: storage, log: log}
}

func (s *TeamServiceImpl) CreateTeam(ctx context.Context, name string) (int64, error) {
	if name == "" {
		return 0, errors.New("name required")
	}
	return s.storage.AddTeam(ctx, name)
}

func (s *TeamServiceImpl) GetTeam(ctx context.Context, id int64) (core.Team, error) {
	team, err := s.storage.GetTeamByID(ctx, id)
	if err != nil {
		return core.Team{}, err
	}
	return team, nil
}

func (s *TeamServiceImpl) GetTeamMembers(ctx context.Context, teamId int64) ([]core.User, error) {
	return s.storage.GetTeamMembers(ctx, teamId)
}
