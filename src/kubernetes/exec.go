package kubernetes

import "strings"

// BuildKubectlArgs appends selected context/namespace to a kubectl argument list.
// Use this for all kubectl invocations from this CLI (except context/namespace commands).
func BuildKubectlArgs(args []string) ([]string, error) {
    contextName, namespace, err := loadKubeSelection()
    if err != nil {
        return nil, err
    }

    if strings.TrimSpace(contextName) != "" {
        args = append(args, "--context", contextName)
    }
    if strings.TrimSpace(namespace) != "" {
        args = append(args, "--namespace", namespace)
    }

    return args, nil
}
