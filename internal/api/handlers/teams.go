package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"Starostina-elena/pull_req_assign/internal/core"
	TeamService "Starostina-elena/pull_req_assign/internal/service"
)

type TeamMemberRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type CreateTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members"`
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type TeamRespMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamResp struct {
	TeamName string           `json:"team_name"`
	Members  []TeamRespMember `json:"members"`
}

func CreateTeamHandler(log *slog.Logger, teamService TeamService.TeamService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateTeamRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		var members []TeamService.TeamMember
		for _, m := range req.Members {
			members = append(members, TeamService.TeamMember{ID: m.UserID, Name: m.Username, IsActive: m.IsActive})
		}

		team, users, err := teamService.CreateTeamWithMembers(r.Context(), req.TeamName, members)
		if err != nil {
			if err == core.ErrTeamExists {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				var er ErrorResponse
				er.Error.Code = "TEAM_EXISTS"
				er.Error.Message = "team with such name already exists"
				_ = json.NewEncoder(w).Encode(er)
				return
			}
			log.Error("failed to create team", "name", req.TeamName, "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var resp TeamResp
		resp.TeamName = team.Name
		for _, u := range users {
			resp.Members = append(resp.Members, TeamRespMember{UserID: u.ID, Username: u.Name, IsActive: u.IsActive})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]TeamResp{"team": resp})
	}
}

func GetTeamHandler(log *slog.Logger, teamService TeamService.TeamService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			http.Error(w, "invalid team id", http.StatusBadRequest)
			return
		}

		team, err := teamService.GetTeam(r.Context(), id)
		if err != nil {
			log.Error("failed to get team", "id", id, "error", err)
			http.Error(w, "No such team", http.StatusNotFound)
			return
		}

		members, err := teamService.GetTeamMembers(r.Context(), id)
		if err != nil {
			log.Error("failed to get team members", "id", id, "error", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var resp TeamResp
		resp.TeamName = team.Name
		for _, u := range members {
			resp.Members = append(resp.Members, TeamRespMember{UserID: u.ID, Username: u.Name, IsActive: u.IsActive})
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}
