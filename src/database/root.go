package database

import (
    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

// NewCommand builds the database parent command.
func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:     "database",
        Aliases: []string{"db"},
        Short:   "Database configuration commands",
    }

    ui.AddHelpCommand(cmd)
    cmd.AddCommand(newConfigCommand())

    return cmd
}
