package docker

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/docker/docker/api/types/container"
)

type RunCMD struct {
	DockerImage   string   `validate:"required" json:"docker_image"`
	ContainerName string   `validate:"required" json:"container_name"`
	Cmd           []string `json:"cmd"`
	Persist       bool     `json:"persist"`
}

func (cmd RunCMD) Validate(whiteListedDockerImages ...string) error {
	err := validator.New().Struct(cmd)
	if err != nil {
		return err
	}
	if len(whiteListedDockerImages) == 0 {
		return nil
	}

	for _, allowed := range whiteListedDockerImages {
		if allowed == cmd.DockerImage {
			return nil
		}
	}
	return fmt.Errorf("image not allowed: %s", cmd.DockerImage)
}

func (cmd RunCMD) ContainerConfig() *container.Config {
	return &container.Config{
		Image: cmd.DockerImage,
		Cmd:   cmd.Cmd,
	}
}
