package service

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"Starostina-elena/pull_req_assign/internal/storage"
	"context"
	"errors"
	"log/slog"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, title string, authorId int64, isOpened bool) (int64, error)
	GetPullRequest(ctx context.Context, id int64) (core.PullRequest, error)
}

type PullRequestServiceImpl struct {
	storage *storage.DB
	log     *slog.Logger
}

func NewPullRequestService(log *slog.Logger, storage *storage.DB) PullRequestService {
	return &PullRequestServiceImpl{storage: storage, log: log}
}

func (s *PullRequestServiceImpl) CreatePullRequest(ctx context.Context, title string, authorId int64, isOpened bool) (int64, error) {
	if title == "" {
		return 0, errors.New("title required")
	}
	return s.storage.AddPullRequest(ctx, title, authorId, isOpened, nil, nil)
}

func (s *PullRequestServiceImpl) GetPullRequest(ctx context.Context, id int64) (core.PullRequest, error) {
	pr, err := s.storage.GetPullRequestByID(ctx, id)
	if err != nil {
		return core.PullRequest{}, err
	}
	return pr, nil
}
