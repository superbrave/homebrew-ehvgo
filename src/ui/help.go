package ui

import "github.com/spf13/cobra"

// AddHelpCommand attaches a "help" subcommand that shows help for the parent command.
func AddHelpCommand(target *cobra.Command) {
    if target == nil {
        return
    }

    helpCmd := &cobra.Command{
        Use:   "help",
        Short: "Show help for this command",
        Args:  cobra.ArbitraryArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            parent := cmd.Parent()
            if parent == nil {
                return cmd.Help()
            }
            return parent.Help()
        },
    }

    target.SetHelpCommand(helpCmd)
}
