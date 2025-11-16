package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	TeamService "Starostina-elena/pull_req_assign/internal/service"
)

type CreateTeamRequest struct {
	Name string `json:"name"`
}

type CreateTeamResponse struct {
	ID int64 `json:"id"`
}

func CreateTeamHandler(log *slog.Logger, teamService TeamService.TeamService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTeamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		id, err := teamService.CreateTeam(r.Context(), req.Name)
		if err != nil {
			log.Error("failed to create team", "error", err)
			http.Error(w, "Error while creating team", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp := CreateTeamResponse{ID: id}
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func GetTeamHandler(log *slog.Logger, teamService TeamService.TeamService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, err := teamService.GetTeam(r.Context(), id)
		if err != nil {
			log.Error("failed to get team", "id", id, "error", err)
			http.Error(w, "No such team", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(u)
	}
}
