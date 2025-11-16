package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	PullRequestService "Starostina-elena/pull_req_assign/internal/service"
)

type CreatePullRequestRequest struct {
	Title    string `json:"title"`
	AuthorID int64  `json:"author_id"`
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

func GetPullRequestHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, err := pullRequestService.GetPullRequest(r.Context(), id)
		if err != nil {
			log.Error("failed to get pull request", "id", id, "error", err)
			http.Error(w, "No such pull request", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(u)
	}
}
