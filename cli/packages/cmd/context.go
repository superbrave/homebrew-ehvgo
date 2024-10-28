package cmd

import (
	"archive/tar"
	"context"
	"ehvg/packages/util"
	"errors"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
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

		color.New(color.FgHiYellow).Println("Starting Infisical CLI..")

		ctx := context.Background()
		projectContext := args[0]

		docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		util.HandleError(err, true)

		defer docker.Close()

		pull, err := docker.ImagePull(ctx, util.INFISICAL_CLI_IMAGE, image.PullOptions{
			RegistryAuth: util.GetDockerAuth("ghcr.io"),
		})
		util.HandleError(err, true)

		defer pull.Close()

		io.Copy(io.Discard, pull)

		infisicalClientAuth, err := util.GetInfisicalClientAuth(projectContext)
		util.HandleError(err, true)

		containerConfig := &container.Config{
			Image: util.GetInfisicalmage("latest"),
			Cmd:   []string{"sleep", "20"},
			Tty:   false,
		}

		hostConfig := &container.HostConfig{
			NetworkMode: "bridge",
			AutoRemove:  true,
		}

		c, err := docker.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
		util.HandleError(err, true)

		if err := docker.ContainerStart(ctx, c.ID, container.StartOptions{}); err != nil {
			util.PrintError(err, true)
		}

		cec, err := docker.ContainerExecCreate(ctx, c.ID, container.ExecOptions{
			Env: []string{
				"INFISICAL_UNIVERSAL_AUTH_CLIENT_ID=" + infisicalClientAuth.ClientId,
				"INFISICAL_UNIVERSAL_AUTH_CLIENT_SECRET=" + infisicalClientAuth.ClientSecret,
				"INFISICAL_PROJECT_ID=" + infisicalClientAuth.ProjectId,
			},
			Cmd:          []string{"bash", "-c", "infisical export -edev &> /var/www/.env &"},
			AttachStdin:  true,
			AttachStdout: true,
		})
		util.HandleError(err, true)

		if err := docker.ContainerExecStart(ctx, cec.ID, container.ExecStartOptions{}); err != nil {
			util.HandleError(err, true)
		}

		cei, err := docker.ContainerExecInspect(ctx, cec.ID)
		util.HandleError(err, true)

		if cei.ExitCode == 0 {
			if _, err := os.Stat(".env"); err == nil {
				color.New(color.FgHiYellow).Println(".env file already exists, deleting..")
				os.Remove(".env")
			}

			color.New(color.FgHiYellow).Println("Creating new .env file..")

			srcFile, _, err := docker.CopyFromContainer(context.Background(), c.ID, "/var/www/.env")
			util.HandleError(err, true)

			defer srcFile.Close()

			tr := tar.NewReader(srcFile)

			for {
				header, err := tr.Next()

				if header != nil {
					f, err := os.OpenFile(util.GetCwdForFile(".env"), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
					util.HandleError(err, false)

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

		docker.ContainerStop(ctx, c.ID, container.StopOptions{Signal: "SIGKILL"})
	}
}

func init() {
	contextCmd.AddCommand(setContextCmd)
	rootCmd.AddCommand(contextCmd)
}
