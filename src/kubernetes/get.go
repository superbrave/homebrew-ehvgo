package kubernetes

import (
    "errors"
    "os"
    "os/exec"

    "ehvgo/src/ui"

    "github.com/spf13/cobra"
)

func newGetCommand() *cobra.Command {
    var namespaceOverride string
    var allNamespaces bool

    cmd := &cobra.Command{
        Use:   "get",
        Short: "Get Kubernetes resources",
        Args:  cobra.MinimumNArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            if allNamespaces && namespaceOverride != "" {
                return errors.New("cannot use --all and --namespace together")
            }

            kubectlArgs := append([]string{"get"}, args...)
            resolvedArgs, err := BuildKubectlArgsWithOptions(kubectlArgs, namespaceOverride, allNamespaces)
            if err != nil {
                return err
            }

            execCmd := exec.Command("kubectl", resolvedArgs...)
            execCmd.Stdin = cmd.InOrStdin()

            stopSpinner := ui.StartSpinner(os.Stderr, "Running kubectl get")
            execCmd.Stdout = ui.WrapWriterOnFirstWrite(cmd.OutOrStdout(), stopSpinner)
            execCmd.Stderr = ui.WrapWriterOnFirstWrite(cmd.ErrOrStderr(), stopSpinner)

            err = execCmd.Run()
            stopSpinner()
            if err != nil {
                if errors.Is(err, exec.ErrNotFound) {
                    return errors.New("kubectl not found in PATH")
                }
                return err
            }

            return nil
        },
    }

    cmd.Flags().StringVarP(&namespaceOverride, "namespace", "n", "", "Kubernetes namespace")
    cmd.Flags().BoolVar(&allNamespaces, "all", false, "Use all namespaces")

    ui.AddHelpCommand(cmd)
    return cmd
}
