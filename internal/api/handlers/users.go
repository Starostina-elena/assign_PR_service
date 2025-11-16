package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	userService "Starostina-elena/pull_req_assign/internal/service"
)

type CreateUserRequest struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

func CreateUserHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		err := userService.CreateUser(r.Context(), req.Id, req.Name, req.IsActive)
		if err != nil {
			log.Error("failed to create user", "error", err)
			http.Error(w, "Error while creating user", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func GetUserHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := userService.GetUser(r.Context(), r.PathValue("id"))
		if err != nil {
			log.Error("failed to get user", "id", r.PathValue("id"), "error", err)
			http.Error(w, "No such user", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(u)
	}
}

func SetUserTeamHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teamId, err := strconv.ParseInt(r.PathValue("team_id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid team id", http.StatusBadRequest)
			return
		}
		err = userService.SetTeamToUser(r.Context(), r.PathValue("user_id"), teamId)
		if err != nil {
			log.Error("failed to set user team", "user_id", r.PathValue("user_id"), "team_id", teamId, "error", err)
			http.Error(w, "Error while setting user team", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func ExpellUserfromTeamHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := userService.ExpellUserFromTeam(r.Context(), r.PathValue("user_id"))
		if err != nil {
			log.Error("failed to expel user from team", "user_id", r.PathValue("user_id"), "error", err)
			http.Error(w, "Error while expelling user from team", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

type ActivateUserRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func SetUserIsActiveHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ActivateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		err := userService.SetUserIsActive(r.Context(), req.UserID, req.IsActive)
		log.Info("Start setting user is active", "user id", req.UserID, "is active", req.IsActive)
		if err != nil {
			log.Error("failed to set user is active", "user_id", req.UserID, "error", err)
			http.Error(w, "Error while setting user is active", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type GetPullRequestsAssignedResponse struct {
	UserId   string                        `json:"user_id"`
	OpenedPr []GetPullRequestShortResponse `json:"pull_requests"`
}

type GetPullRequestShortResponse struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	AuthorID string `json:"author_id"`
	Status   string `json:"is_opened"`
}

func GetPullRequestsAssignedHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prs, err := userService.GetPullRequestAssigned(r.Context(), r.PathValue("user_id"))
		resp := GetPullRequestsAssignedResponse{}
		for _, pr := range prs {
			prResp := GetPullRequestShortResponse{
				ID:       pr.ID,
				Title:    pr.Title,
				AuthorID: pr.AuthorID,
			}
			if pr.IsOpened {
				prResp.Status = "OPEN"
			} else {
				prResp.Status = "MERGED"
			}
			resp.OpenedPr = append(resp.OpenedPr, prResp)
		}
		resp.UserId = r.PathValue("user_id")
		if err != nil {
			log.Error("failed to get pull requests assigned to user", "user_id", r.PathValue("user_id"), "error", err)
			http.Error(w, "Error while getting pull requests assigned to user", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
