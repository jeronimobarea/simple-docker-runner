package docker

import (
	"context"
)

//go:generate mockgen -destination=../test/dockertest/mocks/mock_service.go -package=mocks . Service
type Service interface {
	Run(ctx context.Context, cmd RunCMD) (*RunnerResponse, error)
}

type service struct {
	runner                  DockerRunner
	whiteListedDockerImages []string
}

func NewService(runner DockerRunner, whiteListedDockerImages ...string) *service {
	return &service{runner, whiteListedDockerImages}
}

func (svc service) Run(ctx context.Context, cmd RunCMD) (*RunnerResponse, error) {
	if err := cmd.Validate(svc.whiteListedDockerImages...); err != nil {
		return nil, err
	}

	return svc.runner.Run(
		ctx,
		cmd.DockerImage,
		cmd.ContainerName,
		cmd.Cmd,
		cmd.ContainerConfig(),
		cmd.Persist,
	)
}
