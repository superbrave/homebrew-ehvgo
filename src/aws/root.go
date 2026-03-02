package aws

import (
    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

// NewCommand builds the aws parent command.
func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "aws",
        Short: "AWS related commands",
    }

    ui.AddHelpCommand(cmd)
    cmd.AddCommand(newLoginCommand())

    return cmd
}
