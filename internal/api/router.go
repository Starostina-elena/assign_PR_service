package api

import (
	"log/slog"
	"net/http"

	"Starostina-elena/pull_req_assign/internal/api/handlers"
	"Starostina-elena/pull_req_assign/internal/service"
)

func NewRouter(logger *slog.Logger, userService service.UserService) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /users", handlers.CreateUserHandler(logger, userService))
	mux.Handle("GET /users/{id}", handlers.GetUserHandler(logger, userService))

	return mux
}
