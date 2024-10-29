package infisical

import (
	"fmt"

	"github.com/spf13/cobra"
)

var env string
var service string
var application string

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Commands to manage the contexts (environments) of your current Docker Compose project",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var setContextCmd = &cobra.Command{
	Use:    "set",
	Short:  "Change the current context of your Docker Compose project",
	PreRun: PreRunSetContext,
	Run:    SetContext,
}

func PreRunSetContext(cmd *cobra.Command, args []string) {

}

func SetContext(cmd *cobra.Command, args []string) {
	env, _ := cmd.Flags().GetString("env")
	//application, _ := cmd.Flags().GetString("application")
	//project, _ := cmd.Flags().GetString("project")

	//client := util.
	fmt.Println(env)
}

func init() {
	setContextCmd.Flags().StringVarP(&env, "env", "e", "dev", "The environment to retrieve the variables for")
	setContextCmd.Flags().StringVarP(&service, "service", "", "", "The service to retrieve the variables for")
	setContextCmd.Flags().StringVarP(&application, "application", "", "", "The application to retrieve the variables for")
	contextCmd.AddCommand(setContextCmd)
	rootCmd.AddCommand(contextCmd)
}
