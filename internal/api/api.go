package api

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apiHandlers "github.com/jeronimobarea/simple-docker-runner/internal/api/handlers"
	"github.com/jeronimobarea/simple-docker-runner/internal/docker"
	dockerHandlers "github.com/jeronimobarea/simple-docker-runner/internal/docker/handlers"
	"github.com/jeronimobarea/simple-docker-runner/internal/pkg/env"
)

var serverPort string

func init() {
	serverPort = env.GetEnvWithFallback("SERVER_PORT", ":3000")
}

func Run(_ context.Context) {
	var dockerCli *client.Client
	{
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}
		dockerCli = cli
	}

	var dockerSvc docker.Service
	{
		dockerRunner := docker.NewDockerRunner(dockerCli)
		dockerSvc = docker.NewService(dockerRunner)
	}

	var router *chi.Mux
	{
		router = chi.NewRouter()
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.RequestID)

		// Register routes
		apiHandlers.RegisterRoutes(router)
		dockerHandlers.RegisterRoutes(router, dockerSvc)
	}

	http.ListenAndServe(serverPort, router)
}
