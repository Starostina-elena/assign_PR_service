package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	userService "Starostina-elena/pull_req_assign/internal/service"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type CreateUserResponse struct {
	ID int64 `json:"id"`
}

func CreateUserHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id, err := userService.CreateUser(r.Context(), req.Name, req.IsActive)
		if err != nil {
			log.Error("failed to create user", "error", err)
			http.Error(w, "Error while creating user", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp := CreateUserResponse{ID: id}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func GetUserHandler(log *slog.Logger, userService userService.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, err := userService.GetUser(r.Context(), id)
		if err != nil {
			log.Error("failed to get user", "id", id, "error", err)
			http.Error(w, "Error while getting user", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(u)
	}
}
