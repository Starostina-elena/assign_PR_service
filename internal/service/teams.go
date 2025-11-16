package service

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"Starostina-elena/pull_req_assign/internal/storage"
	"context"
	"errors"
	"log/slog"
	"strings"
)

type TeamService interface {
	CreateTeam(ctx context.Context, name string) (int64, error)
	CreateTeamWithMembers(ctx context.Context, name string, members []TeamMember) (core.Team, []core.User, error)
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

type TeamMember struct {
	ID       string
	Name     string
	IsActive bool
}

func (s *TeamServiceImpl) CreateTeamWithMembers(ctx context.Context, name string, members []TeamMember) (core.Team, []core.User, error) {
	if name == "" {
		return core.Team{}, nil, errors.New("name required")
	}

	var coreMembers []core.User
	for _, m := range members {
		coreMembers = append(coreMembers, core.User{ID: m.ID, Name: m.Name, IsActive: m.IsActive})
	}

	team, users, err := s.storage.AddTeamWithMembers(ctx, name, coreMembers)
	if err != nil {
		if err.Error() != "" && (containsIgnoreCase(err.Error(), "duplicate") || containsIgnoreCase(err.Error(), "unique")) {
			return core.Team{}, nil, core.ErrTeamExists
		}
		return core.Team{}, nil, err
	}

	return team, users, nil
}

func containsIgnoreCase(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
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
