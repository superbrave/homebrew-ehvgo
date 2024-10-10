package util

import (
	"os/exec"
)

const DOCKER_COMPOSE_FILENAME = "docker-compose.yaml"

func HasDocker() bool {
	_, err := exec.LookPath("docker")

	if err != nil {
		PrintError(err.Error())

		return false
	}

	return true
}
