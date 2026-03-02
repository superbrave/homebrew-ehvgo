package kubernetes

import (
    "context"
    "errors"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "time"

    "ehvgo/src/ui"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func newSetNamespaceCommand() *cobra.Command {
    var namespaceFlag string

    cmd := &cobra.Command{
        Use:   "set-namespace",
        Short: "Select a Kubernetes namespace",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            namespace := strings.TrimSpace(namespaceFlag)
            if namespace == "" {
                cfg, err := readConfig()
                if err != nil {
                    return err
                }

                var namespaces []string
                err = ui.RunWithSpinner(os.Stderr, "Fetching namespaces", func() error {
                    var listErr error
                    namespaces, listErr = listNamespaces()
                    return listErr
                })
                if err != nil {
                    return err
                }

                selectPrompt := promptui.Select{
                    Label:  "Select Kubernetes namespace",
                    Items:  namespaces,
                    Size:   10,
                    Stdout: bellSkipper{},
                }

                _, result, err := selectPrompt.Run()
                if err != nil {
                    if errors.Is(err, promptui.ErrAbort) {
                        return nil
                    }
                    return err
                }

                namespace = strings.TrimSpace(result)
                if namespace == "" {
                    namespace = strings.TrimSpace(cfg.KubeNamespace)
                }
            }

            if namespace == "" {
                return errors.New("namespace is required")
            }

            cfg, err := readConfig()
            if err != nil {
                return err
            }
            cfg.KubeNamespace = namespace

            if err := writeConfig(cfg); err != nil {
                return err
            }

            _, err = cmd.OutOrStdout().Write([]byte("Selected namespace: " + namespace + "\n"))
            return err
        },
    }

    cmd.Flags().StringVar(&namespaceFlag, "namespace", "", "Kubernetes namespace")

    ui.AddHelpCommand(cmd)
    return cmd
}

func listNamespaces() ([]string, error) {
    contextName, _, err := loadKubeSelection()
    if err != nil {
        return nil, err
    }
    if strings.TrimSpace(contextName) == "" {
        return nil, errors.New("no context set; run 'ehvg k8s set-context'")
    }

    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    kubeconfig := filepath.Join(home, ".kube", "config")
    overrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
    loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
    clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)

    restConfig, err := clientConfig.ClientConfig()
    if err != nil {
        return nil, err
    }

    clientset, err := kubernetes.NewForConfig(restConfig)
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    list, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    if err != nil {
        return nil, err
    }

    if len(list.Items) == 0 {
        return nil, errors.New("no namespaces found in cluster")
    }

    namespaces := make([]string, 0, len(list.Items))
    for _, item := range list.Items {
        name := strings.TrimSpace(item.Name)
        if name != "" {
            namespaces = append(namespaces, name)
        }
    }
    sort.Strings(namespaces)

    return namespaces, nil
}
