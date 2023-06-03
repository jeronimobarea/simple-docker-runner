package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
)

func TestDockerService(t *testing.T) {
	dockerCli, err := client.NewClientWithOpts(client.FromEnv)
	require.NoError(t, err)

	var (
		ctx           = context.Background()
		allowedImages = []string{"hello-world"}

		dockerRunner = NewDockerRunner(dockerCli)
		dockerSvc    = NewService(dockerRunner, allowedImages...)
	)

	t.Run("happy path", func(t *testing.T) {
		cmd := RunCMD{
			DockerImage:   allowedImages[0],
			ContainerName: "docker-svc-test",
			Persist:       false,
		}

		expected := `{
			"output": [
				"",
				"Hello from Docker!",
				"This message shows that your installation appears to be working correctly.",
				"",
				"To generate this message, Docker took the following steps:",
				" 1. The Docker client contacted the Docker daemon.",
				" 2. The Docker daemon pulled the \"hello-world\" image from the Docker Hub.",
				"    (amd64)",
				" 3. The Docker daemon created a new container from that image which runs the",
				"    executable that produces the output you are currently reading.",
				" 4. The Docker daemon streamed that output to the Docker client, which sent it",
				"    to your terminal.",
				"",
				"To try something more ambitious, you can run an Ubuntu container with:",
				" $ docker run -it ubuntu bash",
				"",
				"Share images, automate workflows, and more with a free Docker ID:",
				" https://hub.docker.com/",
				"",
				"For more examples and ideas, visit:",
				" https://docs.docker.com/get-started/",
				""
			]
		}`

		res, err := dockerSvc.Run(ctx, cmd)
		require.NoError(t, err)
		require.EqualValues(t, expected, res)
		ensureContainerIsRemoved(t, ctx, dockerCli, cmd.ContainerName)
	})

	t.Run("fail to execute non whitelisted image", func(t *testing.T) {
		cmd := RunCMD{
			DockerImage:   "non-whitelisted-image",
			ContainerName: "non-whitelisted-image-test",
			Persist:       false,
		}

		res, err := dockerSvc.Run(ctx, cmd)
		require.EqualError(t, err, "image not allowed: non-whitelisted-image")
		require.Nil(t, res)
	})
}

func ensureContainerIsRemoved(t *testing.T, ctx context.Context, cli *client.Client, containerName string) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	require.NoError(t, err)
	for _, container := range containers {
		if container.Names[0] == containerName {
			require.Fail(t, "conainer not removed")
		}
	}
}
