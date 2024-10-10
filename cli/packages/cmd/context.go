package cmd

import (
	"ehvg/packages/util"
	"fmt"

	"github.com/spf13/cobra"
)

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
		context, err := cmd.Flags().GetString("context")

		if err != nil {
			util.PrintError(err.Error())
		}

		fmt.Println(context)
	}
}

func init() {
	setContextCmd.Flags().String("context", "", "Name of the context to set")
	contextCmd.AddCommand(setContextCmd)
	rootCmd.AddCommand(contextCmd)
}
