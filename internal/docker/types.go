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

	//HostIP   string `validate:"required" json:"host_ip"`
	//HostPort string `validate:"required" json:"host_port"`

	Persist bool `json:"persist"`
}

func (cmd RunCMD) Validate(allowedDockerImages ...string) error {
	err := validator.New().Struct(cmd)
	if err != nil {
		return err
	}
	if len(allowedDockerImages) == 0 {
		return nil
	}
	for _, allowed := range allowedDockerImages {
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

func (cmd RunCMD) HostConfig() *container.HostConfig {
	return &container.HostConfig{}
}
