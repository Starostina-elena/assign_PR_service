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
	mux.Handle("PUT /users/setIsActive", handlers.SetUserIsActiveHandler(logger, userService))
	mux.Handle("GET /users/{user_id}/getReview", handlers.GetPullRequestsAssignedHandler(logger, userService))

	mux.Handle("POST /team/add", handlers.CreateTeamHandler(logger, teamService))
	mux.Handle("GET /team/{id}/get", handlers.GetTeamHandler(logger, teamService))

	mux.Handle("POST /pullRequest/create", handlers.CreatePullRequestHandler(logger, pullRequestService))
	mux.Handle("GET /pullRequest/{id}", handlers.GetPullRequestHandler(logger, pullRequestService))
	mux.Handle("PUT /pullRequest/{pull_request_id}/reviewer/{num_reviewer}", handlers.ChangeReviewerHandler(logger, pullRequestService))
	mux.Handle("PUT /pullRequest/merge", handlers.MergePullRequestHandler(logger, pullRequestService))
	mux.Handle("POST /pullRequest/reassign", handlers.ReassignPullRequestHandler(logger, pullRequestService))

	return mux
}
