package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	apiHandlers "github.com/jeronimobarea/simple-docker-runner/internal/api/handlers"
	"github.com/jeronimobarea/simple-docker-runner/internal/docker"
	dockerHandlers "github.com/jeronimobarea/simple-docker-runner/internal/docker/handlers"
	"github.com/jeronimobarea/simple-docker-runner/internal/pkg/env"
)

var (
	serverPort          string
	allowedDockerImages []string
)

func init() {
	allowedDockerImagesFilePath := os.Getenv("ALLOWED_DOCKER_IMAGES")
	if allowedDockerImagesFilePath != "" {
		if !strings.HasSuffix(allowedDockerImagesFilePath, ".json") {
			panic("file has to be a .json file")
		}
		allowedDockerImagesFile, err := os.Open(allowedDockerImagesFilePath)
		if err != nil {
			panic(err)
		}
		defer allowedDockerImagesFile.Close()

		decoder := json.NewDecoder(allowedDockerImagesFile)
		err = decoder.Decode(&allowedDockerImages)
		if err != nil {
			panic(err)
		}
	}

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
		dockerSvc = docker.NewService(dockerRunner, allowedDockerImages...)
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
