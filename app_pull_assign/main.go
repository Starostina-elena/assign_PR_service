package main

import (
	"Starostina-elena/pull_req_assign/internal"
	"log/slog"
)

func Run() {
	slog.Info("Server started")
	err := internal.Init()
	if err != nil {
		slog.Error("Initialization failed", "error", err)
		return
	}
	slog.Info("Initialization succeeded")
}

func main() {
	Run()
}
