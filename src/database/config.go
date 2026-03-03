package database

import "github.com/spf13/cobra"

func newConfigCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "config",
        Short: "Configure database connections",
    }

    cmd.AddCommand(newConfigAddCommand())
    return cmd
}
