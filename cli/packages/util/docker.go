package util

import (
	"os/exec"
)

type DockerEvent struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func HasDocker() bool {
	_, err := exec.LookPath("docker")

	if err != nil {
		PrintError(err.Error())

		return false
	}

	return true
}

func GetInfisicalmage(version string) string {
	return INFISICAL_CLI_IMAGE + ":" + version
}
