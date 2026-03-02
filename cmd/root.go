package cmd

import (
    "fmt"
    "os"

    "ehvgo/src/aws"
    "ehvgo/src/kubernetes"
    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "ehvg",
    Short: "ehvg is a CLI application",
    Long:  "ehvg is a CLI application built with Cobra.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Fprintln(cmd.OutOrStdout(), "ehvg CLI")
    },
}

func init() {
    ui.AddHelpCommand(rootCmd)
    rootCmd.AddCommand(aws.NewCommand())
    rootCmd.AddCommand(kubernetes.NewCommand())
}

// Execute runs the root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
