package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"Starostina-elena/pull_req_assign/internal/core"
	PullRequestService "Starostina-elena/pull_req_assign/internal/service"
)

type CreatePullRequestRequest struct {
	Title    string `json:"title"`
	AuthorID string `json:"author_id"`
	IsOpened bool   `json:"is_opened"`
}

type CreatePullRequestResponse struct {
	ID int64 `json:"id"`
}

func CreatePullRequestHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePullRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id, err := pullRequestService.CreatePullRequest(r.Context(), req.Title, req.AuthorID, req.IsOpened)
		if err != nil {
			log.Error("failed to create pull request", "error", err)
			http.Error(w, "Error while creating pull request", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp := CreatePullRequestResponse{ID: id}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

type GetPullRequestResponse struct {
	ID        int64   `json:"id"`
	Title     string  `json:"title"`
	AuthorID  string  `json:"author_id"`
	Reviewer1 *string `json:"reviewer_1,omitempty"`
	Reviewer2 *string `json:"reviewer_2,omitempty"`
	Status    string  `json:"is_opened"`
}

func GetPullRequestHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		pr, err := pullRequestService.GetPullRequest(r.Context(), id)
		if err != nil {
			log.Error("failed to get pull request", "id", id, "error", err)
			http.Error(w, "No such pull request", http.StatusNotFound)
			return
		}

		resp := GetPullRequestResponse{
			ID:       pr.ID,
			Title:    pr.Title,
			AuthorID: pr.AuthorID,
		}
		if pr.IsOpened {
			resp.Status = "OPEN"
		} else {
			resp.Status = "MERGED"
		}
		if pr.Reviewer1ID.Valid {
			resp.Reviewer1 = &pr.Reviewer1ID.String
		}
		if pr.Reviewer2ID.Valid {
			resp.Reviewer2 = &pr.Reviewer2ID.String
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func MergePullRequestHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		err = pullRequestService.MergePullRequest(r.Context(), id)
		if err != nil {
			log.Error("failed to merge pull request", "id", id, "error", err)
			http.Error(w, "Error while merging pull request", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ChangeReviewerHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pullReqId, err := strconv.ParseInt(r.PathValue("pull_request_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid pull request id", http.StatusBadRequest)
			return
		}
		numReviwer, err := strconv.ParseInt(r.PathValue("num_reviewer"), 10, 64)
		if err != nil || (numReviwer != 1 && numReviwer != 2) {
			http.Error(w, "invalid reviewer number", http.StatusBadRequest)
			return
		}

		err = pullRequestService.ChangeReviewer(r.Context(), pullReqId, numReviwer)
		if err != nil {
			if errors.Is(err, core.ErrNotEnoughCoworkers) || errors.Is(err, core.PullRequestAlreadyMerged) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			log.Error("failed to change reviewer", "pull_request_id", pullReqId, "num_reviewer", numReviwer, "error", err)
			http.Error(w, "Error while changing reviewer", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
