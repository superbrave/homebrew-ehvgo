package kubernetes

import (
    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

// NewCommand builds the k8s parent command.
func NewCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:     "k8s",
        Aliases: []string{"kubernetes"},
        Short:   "Kubernetes commands",
    }

    ui.AddHelpCommand(cmd)
    cmd.AddCommand(newSetContextCommand())
    cmd.AddCommand(newShowContextCommand())
    cmd.AddCommand(newSetNamespaceCommand())
    cmd.AddCommand(newShowNamespaceCommand())
    cmd.AddCommand(newGetCommand())
    cmd.AddCommand(newExecCommand())

    return cmd
}
