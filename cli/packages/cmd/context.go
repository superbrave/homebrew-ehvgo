package cmd

import (
	"archive/tar"
	"context"
	"ehvg/packages/util"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Commands to manage the contexts (environments) of your current Docker Compose project",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var setContextCmd = &cobra.Command{
	Use:       "set",
	Short:     "Change the current context of your Docker Compose project",
	Run:       SetContext,
	ValidArgs: []string{"dok", "seeme"},
}

func SetContext(cmd *cobra.Command, args []string) {
	if util.HasDocker() {
		if len(args) != 1 {
			util.PrintError(errors.New("missing context"), true)
		}

		ctx := context.Background()
		projectContext := args[0]

		docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			util.PrintError(err, true)
		}

		defer docker.Close()

		pull := util.PullImage(docker, util.GetInfisicalmage("latest"))

		io.Copy(os.Stdout, pull)

		infisicalClientAuth, err := util.GetInfisicalClientAuth(projectContext)
		if err != nil {
			util.PrintError(err, true)
		}

		containerConfig := &container.Config{
			Image: util.GetInfisicalmage("latest"),
			Env: []string{
				"INFISICAL_UNIVERSAL_AUTH_CLIENT_ID=" + infisicalClientAuth.ClientId,
				"INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET=" + infisicalClientAuth.ClientSecret,
				"INFISICAL_PROJECT_ID=" + infisicalClientAuth.ProjectId,
			},
			Cmd: []string{
				"/bin/bash",
				"-c",
				"touch /var/www/.env && infisical login",
			},
			Shell:        []string{"/bin/bash"},
			Tty:          true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
		}

		c, err := docker.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
		if err != nil {
			util.PrintError(err, true)
		}

		if err := docker.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
			util.PrintError(err, true)
		}

		defer util.StopAndRemoveContainer(docker, ctx, c.ID)

		out, err := docker.ContainerLogs(ctx, c.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
		if err != nil {
			fmt.Printf("Error getting logs: %v\n", err)
			return
		}
		defer out.Close()

		// Print the logs
		fmt.Println("Container logs:")
		_, err = io.Copy(os.Stdout, out)
		if err != nil {
			fmt.Printf("Error reading logs: %v\n", err)
		}

		srcFile, _, err := docker.CopyFromContainer(context.Background(), c.ID, "/var/www/.env")
		if err != nil {
			util.PrintError(err, true)
		}

		defer srcFile.Close()

		tr := tar.NewReader(srcFile)

		for {
			header, err := tr.Next()

			if header != nil {
				f, err := os.OpenFile(util.GetCwdForFile(".env"), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					util.PrintError(err, false)
				}

				if _, err := io.Copy(f, tr); err != nil {
					util.PrintError(err, false)
				}

				f.Close()
			}

			if err != nil {
				break
			}
		}
	}

}

func init() {
	contextCmd.AddCommand(setContextCmd)
	rootCmd.AddCommand(contextCmd)
}
