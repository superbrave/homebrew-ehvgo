package cmd

import (
	"context"
	"ehvg/packages/util"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var event *util.DockerEvent

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Commands to manage the contexts (environments) of your current Docker Compose project",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var setContextCmd = &cobra.Command{
	Use:   "set",
	Short: "Change the current context of your Docker Compose project",
	Run:   SetContext,
}

func SetContext(cmd *cobra.Command, args []string) {
	if util.HasDocker() {
		projectContext, err := cmd.Flags().GetString("context")

		if err != nil {
			util.PrintError(err.Error())
		}

		docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

		if err != nil {
			util.PrintError("Unable to create Docker client")
		}

		pullResult, err := docker.ImagePull(context.Background(), util.GetInfisicalmage("latest"), image.PullOptions{})

		if err != nil {
			pullResult.Close()
			util.PrintError(err.Error())
		}

		events := json.NewDecoder(pullResult)

		var event util.DockerEvent

		for {
			if err := events.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
			}

			fmt.Printf("EVENT: %+v\n", event)
		}

		fmt.Println(projectContext)
	}
}

func init() {
	setContextCmd.Flags().String("context", "", "Name of the context to set")
	contextCmd.AddCommand(setContextCmd)
	rootCmd.AddCommand(contextCmd)
}
