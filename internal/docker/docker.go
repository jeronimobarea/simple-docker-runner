package docker

import (
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var ErrContainerStop = errors.New("error stoping container")

//go:generate mockgen -destination=../test/dockertest/mocks/mock_runner.go -package=mocks . DockerRunner
type DockerRunner interface {
	Run(ctx context.Context, dockerImage, containerName string, cmd []string, config *container.Config, hostConfig *container.HostConfig, persist bool) ([]byte, error)
}

type dockerRunner struct {
	cli *client.Client
}

func NewDockerRunner(cli *client.Client) *dockerRunner {
	return &dockerRunner{cli}
}

func (d dockerRunner) Run(
	ctx context.Context,
	dockerImage,
	containerName string,
	cmd []string,
	config *container.Config,
	hostConfig *container.HostConfig,
	persist bool,
) ([]byte, error) {
	containerID, err := d.checkContainerExists(ctx, containerName)
	if err != nil {
		return nil, err
	}
	if containerID == "" {
		containerID, err = d.createContainer(ctx, dockerImage, config, hostConfig, containerName)
		if err != nil {
			return nil, err
		}
	}
	if !persist {
		defer d.deleteContainer(ctx, containerID)
	}
	defer d.stopContainer(ctx, containerID)

	err = d.startContainer(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("error starting container %s: %w", containerID, err)
	}

	statusCh, errCh := d.cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	return d.getContainerOutput(ctx, containerID)
}

func (d dockerRunner) createContainer(
	ctx context.Context,
	image string,
	config *container.Config,
	hostConfig *container.HostConfig,
	containerName string,
) (string, error) {
	reader, err := d.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	reader.Close()

	resp, err := d.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (d dockerRunner) startContainer(ctx context.Context, containerID string) error {
	return d.cli.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (d dockerRunner) stopContainer(ctx context.Context, containerID string) error {
	err := d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return ErrContainerStop
	}
	return nil
}

func (d dockerRunner) getContainerOutput(ctx context.Context, containerID string) ([]byte, error) {
	logs, err := d.cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	var output []byte
	logs.Read(output)
	return output, nil
}

func (d dockerRunner) checkContainerExists(ctx context.Context, containerName string) (string, error) {
	containers, err := d.cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}
	for _, container := range containers {
		if container.Names[0] == containerName {
			return container.ID, nil
		}
	}
	return "", nil
}

func (d dockerRunner) deleteContainer(ctx context.Context, containerID string) error {
	f := func() error {
		return d.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	}
	return doWithRetry(ctx, 3, f)
}

func doWithRetry(ctx context.Context, attempts uint8, f func() error) error {
	if attempts == 0 {
		return nil
	}

	err := f()
	if err != nil {
		return doWithRetry(ctx, attempts-1, f)
	}
	return nil
}
