package util

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
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

type DockerConfig struct {
	Auths map[string]Auths `json:"auths"`
}

type Auths struct {
	Auth string `json:"auth"`
}

func HasDocker() bool {
	_, err := exec.LookPath("docker")

	if err != nil {
		PrintError(err, true)

		return false
	}

	return true
}

// Get authentication info of the current user.
func GetDockerAuth(provider string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		PrintError(err, true)
	}

	configPath := filepath.Join(homeDir, ".docker", "config.json")

	data, err := os.ReadFile(configPath)
	if err != nil {
		PrintError(err, true)
	}

	var dockerConfig DockerConfig

	err = json.Unmarshal(data, &dockerConfig)
	if err != nil {
		PrintError(err, true)
	}

	authString, err := base64.StdEncoding.DecodeString(dockerConfig.Auths[provider].Auth)
	if err != nil {
		PrintError(err, true)
	}

	credentials := strings.Split(string(authString), ":")
	if len(credentials) != 2 {
		PrintError(fmt.Errorf("invalid credentials for registry %v", provider), true)
	}

	authConfig := types.AuthConfig{
		Username:      credentials[0],
		Password:      credentials[1],
		ServerAddress: provider,
	}

	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		PrintError(fmt.Errorf("error encoding credentials: %v", err), true)

	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)

	return authStr
}

func PullImage(docker *client.Client, imageName string) io.ReadCloser {
	pull, err := docker.ImagePull(context.Background(), imageName, image.PullOptions{
		RegistryAuth: GetDockerAuth("ghcr.io"),
	})

	if err != nil {
		PrintError(err, true)
	}

	defer pull.Close()

	events := json.NewDecoder(pull)

	var event DockerEvent

	for {
		if err := events.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
		}

		fmt.Printf("%+v\n", event.Status)
	}

	return pull
}

func StopAndRemoveContainer(docker *client.Client, ctx context.Context, cID string) {
	if err := docker.ContainerStop(ctx, cID, container.StopOptions{Signal: "SIGKILL"}); err != nil {
		PrintError(err, true)
	}

	statusCh, errCh := docker.ContainerWait(ctx, cID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			PrintError(fmt.Errorf("error waiting for container: %v", err), true)
		}
	case status := <-statusCh:
		if status.StatusCode == 0 {
			if err := docker.ContainerRemove(ctx, cID, container.RemoveOptions{Force: true}); err != nil {
				PrintError(err, true)
			}

			fmt.Print("Container successfully removed")
		}
	}
}
