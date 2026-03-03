package cmd

import (
    "context"
    "errors"
    "fmt"
    "os"

    "ehvgo/src/aws"
    "ehvgo/src/database"
    "ehvgo/src/kubernetes"
    "ehvgo/src/ui"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "ehvgo",
    Short: "ehvgo is a CLI application",
    Long:  "ehvgo is a CLI application built with Cobra.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Fprintln(cmd.OutOrStdout(), "ehvgo CLI")
    },
    SilenceErrors: true,
    SilenceUsage:  true,
}

func init() {
    ui.AddHelpCommand(rootCmd)
    rootCmd.AddCommand(aws.NewCommand())
    rootCmd.AddCommand(database.NewCommand())
    rootCmd.AddCommand(kubernetes.NewCommand())
}

// Execute runs the root command.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        if isSilentExit(err) {
            return
        }
        fmt.Fprintln(os.Stderr, err)
        return
    }
}

func isSilentExit(err error) bool {
    if err == nil {
        return true
    }
    if errors.Is(err, promptui.ErrAbort) {
        return true
    }
    if errors.Is(err, promptui.ErrInterrupt) {
        return true
    }
    if errors.Is(err, promptui.ErrEOF) {
        return true
    }
    if errors.Is(err, context.Canceled) {
        return true
    }
    return false
}
