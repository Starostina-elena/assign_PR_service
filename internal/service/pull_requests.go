package service

import (
	"Starostina-elena/pull_req_assign/internal/core"
	"Starostina-elena/pull_req_assign/internal/storage"
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"time"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, title string, authorId string, isOpened bool) (int64, error)
	GetPullRequest(ctx context.Context, id int64) (core.PullRequest, error)
	MergePullRequest(ctx context.Context, id int64) error
	ChangeReviewer(ctx context.Context, pullReqId int64, numReviwer int64) error
}

type PullRequestServiceImpl struct {
	storage *storage.DB
	log     *slog.Logger
}

func NewPullRequestService(log *slog.Logger, storage *storage.DB) PullRequestService {
	return &PullRequestServiceImpl{storage: storage, log: log}
}

func (s *PullRequestServiceImpl) CreatePullRequest(ctx context.Context, title string, authorId string, isOpened bool) (int64, error) {
	if title == "" {
		return 0, errors.New("title required")
	}

	var reviewer1 *string
	var reviewer2 *string
	if isOpened {
		coworkers, err := s.storage.GetActiveCoworkers(ctx, authorId)
		if err == nil {
			if len(coworkers) == 1 {
				reviewer1 = &(coworkers[0].ID)
			} else if len(coworkers) >= 2 {
				rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
				rnd.Shuffle(len(coworkers), func(i, j int) { coworkers[i], coworkers[j] = coworkers[j], coworkers[i] })

				reviewer1 = &(coworkers[0].ID)
				reviewer2 = &(coworkers[1].ID)
			}
		}
	}

	return s.storage.AddPullRequest(ctx, title, authorId, isOpened, reviewer1, reviewer2)
}

func (s *PullRequestServiceImpl) GetPullRequest(ctx context.Context, id int64) (core.PullRequest, error) {
	pr, err := s.storage.GetPullRequestByID(ctx, id)
	if err != nil {
		return core.PullRequest{}, err
	}
	return pr, nil
}

func (s *PullRequestServiceImpl) MergePullRequest(ctx context.Context, id int64) error {
	return s.storage.MergePullRequest(ctx, id)
}

// numReviwer = 1 or 2
func (s *PullRequestServiceImpl) ChangeReviewer(ctx context.Context, pullReqId int64, numReviwer int64) error {
	pr, err := s.storage.GetPullRequestByID(ctx, pullReqId)
	if err != nil {
		return err
	}
	if !pr.IsOpened {
		return core.PullRequestAlreadyMerged
	}
	coworkers, err := s.storage.GetActiveCoworkers(ctx, pr.AuthorID)
	if err != nil {
		return err
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Shuffle(len(coworkers), func(i, j int) { coworkers[i], coworkers[j] = coworkers[j], coworkers[i] })

	for _, u := range coworkers {
		if pr.Reviewer1ID.Valid && u.ID == pr.Reviewer1ID.String {
			continue
		}
		if pr.Reviewer2ID.Valid && u.ID == pr.Reviewer2ID.String {
			continue
		}

		switch numReviwer {
		case 1:
			return s.storage.ChangeReviewer1(ctx, pullReqId, u.ID)
		case 2:
			return s.storage.ChangeReviewer2(ctx, pullReqId, u.ID)
		}
	}

	switch numReviwer {
	case 1:
		err = s.storage.ResetReviewer1(ctx, pullReqId)
		if err != nil {
			return err
		}
	case 2:
		err = s.storage.ResetReviewer2(ctx, pullReqId)
		if err != nil {
			return err
		}
	}
	return core.ErrNotEnoughCoworkers // achieve this line only if no suitable coworker found in cycle
}
