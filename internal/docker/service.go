package docker

import (
	"context"
)

//go:generate mockgen -destination=../test/dockertest/mocks/mock_service.go -package=mocks . Service
type Service interface {
	Run(ctx context.Context, cmd RunCMD) ([]byte, error)
}

type service struct {
	runner DockerRunner
}

func NewService(runner DockerRunner) *service {
	return &service{runner}
}

func (svc service) Run(ctx context.Context, cmd RunCMD) ([]byte, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	return svc.runner.Run(
		ctx,
		cmd.DockerImage,
		cmd.ContainerName,
		cmd.Cmd,
		cmd.ContainerConfig(),
		cmd.HostConfig(),
		cmd.Persist,
	)
}
