package infisical

import "github.com/spf13/cobra"

func newChangeSecretCommand() *cobra.Command {
  var cmd = &cobra.Command{
    Use: "set-secret",
    Aliases: []string{"set-secret"},
    Args: cobra.MinimumNArgs(1),
  }

  return cmd
}