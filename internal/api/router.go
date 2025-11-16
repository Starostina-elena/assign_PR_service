package api

import (
	"log/slog"
	"net/http"

	"Starostina-elena/pull_req_assign/internal/api/handlers"
	"Starostina-elena/pull_req_assign/internal/service"
)

func NewRouter(logger *slog.Logger, userService service.UserService,
	teamService service.TeamService, pullRequestService service.PullRequestService) http.Handler {

	mux := http.NewServeMux()

	mux.Handle("POST /users", handlers.CreateUserHandler(logger, userService))
	mux.Handle("GET /users/{id}", handlers.GetUserHandler(logger, userService))

	mux.Handle("POST /teams", handlers.CreateTeamHandler(logger, teamService))
	mux.Handle("GET /teams/{id}", handlers.GetTeamHandler(logger, teamService))

	mux.Handle("POST /pull_requests", handlers.CreatePullRequestHandler(logger, pullRequestService))
	mux.Handle("GET /pull_requests/{id}", handlers.GetPullRequestHandler(logger, pullRequestService))

	return mux
}
