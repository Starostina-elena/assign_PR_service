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
	mux.Handle("PUT /users/{user_id}/user_team/{team_id}", handlers.SetUserTeamHandler(logger, userService))
	mux.Handle("PUT /users/user_team/{user_id}", handlers.ExpellUserfromTeamHandler(logger, userService))
	mux.Handle("PUT /users/activate/{user_id}", handlers.ActivateUserHandler(logger, userService))
	mux.Handle("PUT /users/deactivate/{user_id}", handlers.DeactivateUserHandler(logger, userService))
	mux.Handle("GET /users/{user_id}/pull_requests_assigned", handlers.GetPullRequestsAssignedHandler(logger, userService))

	mux.Handle("POST /team/add", handlers.CreateTeamHandler(logger, teamService))
	mux.Handle("GET /team/{id}", handlers.GetTeamHandler(logger, teamService))

	mux.Handle("POST /pull_requests", handlers.CreatePullRequestHandler(logger, pullRequestService))
	mux.Handle("GET /pull_requests/{id}", handlers.GetPullRequestHandler(logger, pullRequestService))
	mux.Handle("PUT /pull_requests/{pull_request_id}/reviewer/{num_reviewer}", handlers.ChangeReviewerHandler(logger, pullRequestService))
	mux.Handle("PUT /pull_requests/{id}/merge", handlers.MergePullRequestHandler(logger, pullRequestService))

	return mux
}
