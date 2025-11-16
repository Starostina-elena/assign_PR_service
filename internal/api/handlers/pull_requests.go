package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"Starostina-elena/pull_req_assign/internal/core"
	PullRequestService "Starostina-elena/pull_req_assign/internal/service"
)

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	PR struct {
		PullRequestID     string   `json:"pull_request_id"`
		PullRequestName   string   `json:"pull_request_name"`
		AuthorID          string   `json:"author_id"`
		Status            string   `json:"status"`
		AssignedReviewers []string `json:"assigned_reviewers"`
		CreatedAt         string   `json:"createdAt,omitempty"`
		MergedAt          string   `json:"mergedAt,omitempty"`
	} `json:"pr"`
}

func CreatePullRequestHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePullRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id, err := pullRequestService.CreatePullRequest(r.Context(), req.PullRequestName, req.AuthorID, true)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "foreign") || strings.Contains(strings.ToLower(err.Error()), "referen") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				_ = json.NewEncoder(w).Encode(map[string]map[string]string{"error": {"code": "NOT_FOUND", "message": "author or team not found"}})
				return
			}
			if strings.Contains(strings.ToLower(err.Error()), "unique") || strings.Contains(strings.ToLower(err.Error()), "duplicate") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(map[string]map[string]string{"error": {"code": "PR_EXISTS", "message": "PR id already exists"}})
				return
			}
			log.Error("failed to create pull request", "error", err)
			http.Error(w, "Error while creating pull request", http.StatusInternalServerError)
			return
		}
		pr, err := pullRequestService.GetPullRequest(r.Context(), id)
		if err != nil {
			log.Error("failed to load created pull request", "id", id, "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var assigned []string
		if pr.Reviewer1ID.Valid {
			assigned = append(assigned, pr.Reviewer1ID.String)
		}
		if pr.Reviewer2ID.Valid {
			assigned = append(assigned, pr.Reviewer2ID.String)
		}

		var resp CreatePullRequestResponse
		resp.PR.PullRequestID = req.PullRequestID
		resp.PR.PullRequestName = pr.Title
		resp.PR.AuthorID = pr.AuthorID
		if pr.IsOpened {
			resp.PR.Status = "OPEN"
		} else {
			resp.PR.Status = "MERGED"
		}
		resp.PR.AssignedReviewers = assigned
		resp.PR.CreatedAt = pr.CreatedAt
		resp.PR.MergedAt = pr.MergedAt

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
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

type MergePullRequestRequest struct {
	PullRequestID int64 `json:"pull_request_id"`
}

func MergePullRequestHandler(log *slog.Logger, pullRequestService PullRequestService.PullRequestService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body MergePullRequestRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id := body.PullRequestID

		if err := pullRequestService.MergePullRequest(r.Context(), id); err != nil {
			log.Error("failed to merge pull request", "id", id, "error", err)
			http.Error(w, "Error while merging pull request", http.StatusNotFound)
			return
		}

		pr, err := pullRequestService.GetPullRequest(r.Context(), id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]map[string]string{"error": {"code": "NOT_FOUND", "message": "pull request not found"}})
			return
		}

		var assigned []string
		if pr.Reviewer1ID.Valid {
			assigned = append(assigned, pr.Reviewer1ID.String)
		}
		if pr.Reviewer2ID.Valid {
			assigned = append(assigned, pr.Reviewer2ID.String)
		}

		var resp CreatePullRequestResponse
		resp.PR.PullRequestID = strconv.FormatInt(body.PullRequestID, 10)
		resp.PR.PullRequestName = pr.Title
		resp.PR.AuthorID = pr.AuthorID
		resp.PR.Status = "MERGED"
		resp.PR.AssignedReviewers = assigned
		resp.PR.MergedAt = pr.MergedAt

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
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
