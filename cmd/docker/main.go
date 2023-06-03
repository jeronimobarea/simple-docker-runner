package main

import (
	"context"

	"github.com/jeronimobarea/simple-docker-runner/internal/api"
)

func main() {
	api.Run(context.Background())
}
