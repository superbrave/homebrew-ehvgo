package infisical

import (
	"github.com/spf13/cobra"
)

var infisicalCommand = &cobra.Command {
  Use: "infisical",
  Short: "Infisical Configuration Management Tool",
  Run: func(cmd *cobra.Command, args []string) {},
}

func Execute(rootCmd *cobra.Command) {
  infisicalCommand.AddCommand(
    NewSetEnvironmentCommand(),
  )
  rootCmd.AddCommand(
    infisicalCommand,
  )
}