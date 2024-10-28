package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/docker/cli/cli/config/types"
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
		HandleError(err, false)

		return false
	}

	return true
}

// Get authentication info of the current user.
func GetDockerAuth(provider string) string {
	homeDir, err := os.UserHomeDir()
	HandleError(err, true)

	configPath := filepath.Join(homeDir, ".docker", "config.json")

	data, err := os.ReadFile(configPath)
	HandleError(err, true)

	var dockerConfig DockerConfig

	err = json.Unmarshal(data, &dockerConfig)
	HandleError(err, true)

	authString, err := base64.StdEncoding.DecodeString(dockerConfig.Auths[provider].Auth)
	HandleError(err, true)

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
