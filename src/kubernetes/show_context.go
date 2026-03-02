package kubernetes

import (
    "errors"
    "os"
    "strings"

    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

func newShowContextCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "show-context",
        Short: "Show the selected Kubernetes context",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            var cfg appConfig
            err := ui.RunWithSpinner(os.Stderr, "Loading context", func() error {
                var readErr error
                cfg, readErr = readConfig()
                return readErr
            })
            if err != nil {
                return err
            }

            if strings.TrimSpace(cfg.KubeContext) == "" {
                return errors.New("no context set; run 'ehvg k8s set-context'")
            }

            _, err = cmd.OutOrStdout().Write([]byte(cfg.KubeContext + "\n"))
            return err
        },
    }

    ui.AddHelpCommand(cmd)
    return cmd
}
