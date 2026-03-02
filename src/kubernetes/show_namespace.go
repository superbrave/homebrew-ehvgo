package kubernetes

import (
    "errors"
    "os"
    "strings"

    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

func newShowNamespaceCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "show-namespace",
        Short: "Show the selected Kubernetes namespace",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            var cfg appConfig
            err := ui.RunWithSpinner(os.Stderr, "Loading namespace", func() error {
                var readErr error
                cfg, readErr = readConfig()
                return readErr
            })
            if err != nil {
                return err
            }

            if strings.TrimSpace(cfg.KubeNamespace) == "" {
                return errors.New("no namespace set; run 'ehvg k8s set-namespace'")
            }

            _, err = cmd.OutOrStdout().Write([]byte(cfg.KubeNamespace + "\n"))
            return err
        },
    }

    ui.AddHelpCommand(cmd)
    return cmd
}
