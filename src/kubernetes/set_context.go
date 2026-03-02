package kubernetes

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "sort"
    "strings"

    "ehvgo/src/ui"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
    "k8s.io/client-go/tools/clientcmd"
)

type appConfig struct {
    KubeContext   string `json:"kubeContext"`
    KubeNamespace string `json:"kubeNamespace"`
}

func newSetContextCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "set-context",
        Short: "Select a Kubernetes context",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            var contexts []string
            err := ui.RunWithSpinner(os.Stderr, "Loading contexts", func() error {
                var listErr error
                contexts, listErr = listKubeContexts()
                return listErr
            })
            if err != nil {
                return err
            }
            selected, err := promptForContext(contexts)
            if err != nil {
                return err
            }

            cfg, err := readConfig()
            if err != nil {
                return err
            }
            cfg.KubeContext = selected

            if err := writeConfig(cfg); err != nil {
                return err
            }

            _, err = cmd.OutOrStdout().Write([]byte("Selected context: " + selected + "\n"))
            return err
        },
    }

    ui.AddHelpCommand(cmd)
    return cmd
}

func listKubeContexts() ([]string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    configPath := filepath.Join(home, ".kube", "config")
    config, err := clientcmd.LoadFromFile(configPath)
    if err != nil {
        return nil, err
    }

    if len(config.Contexts) == 0 {
        return nil, errors.New("no contexts found in ~/.kube/config")
    }

    list := make([]string, 0, len(config.Contexts))
    for name := range config.Contexts {
        list = append(list, name)
    }
    sort.Strings(list)

    return list, nil
}

func promptForContext(contexts []string) (string, error) {
    selectPrompt := promptui.Select{
        Label:  "Select Kubernetes context",
        Items:  contexts,
        Size:   10,
        Stdout: bellSkipper{},
    }

    _, result, err := selectPrompt.Run()
    if err != nil {
        return "", err
    }

    return result, nil
}

func readConfig() (appConfig, error) {
    path, err := configPath()
    if err != nil {
        return appConfig{}, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return appConfig{}, nil
        }
        return appConfig{}, err
    }

    var cfg appConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return appConfig{}, err
    }

    cfg.KubeContext = strings.TrimSpace(cfg.KubeContext)
    cfg.KubeNamespace = strings.TrimSpace(cfg.KubeNamespace)

    return cfg, nil
}

func writeConfig(cfg appConfig) error {
    path, err := configPath()
    if err != nil {
        return err
    }

    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0o700); err != nil {
        return err
    }

    payload, err := json.MarshalIndent(cfg, "", "    ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, payload, 0o600)
}

func configPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }

    return filepath.Join(home, ".ehvgo", "config.json"), nil
}

func loadKubeSelection() (string, string, error) {
    cfg, err := readConfig()
    if err != nil {
        return "", "", err
    }

    return cfg.KubeContext, cfg.KubeNamespace, nil
}
