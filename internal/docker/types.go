package docker

import (
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

func (cmd RunCMD) Validate() error {
	return validator.New().Struct(cmd)
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
