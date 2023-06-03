package api

import (
	"context"
	"encoding/json"
	"fmt"
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
	serverPort              string
	whiteListedDockerImages []string
)

func init() {
	whitelistedDockerImagesFilePath := os.Getenv("WHITELISTED_DOCKER_IMAGES")
	if whitelistedDockerImagesFilePath != "" {
		loadWhitelistedDockerImages(whitelistedDockerImagesFilePath)
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
		dockerSvc = docker.NewService(dockerRunner, whiteListedDockerImages...)
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

func loadWhitelistedDockerImages(filePath string) {
	if !strings.HasSuffix(filePath, ".json") {
		panic(fmt.Errorf("file has to be a .json file: %s", filePath))
	}
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&whiteListedDockerImages)
	if err != nil {
		panic(err)
	}
}
