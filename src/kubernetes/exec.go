package kubernetes

import "strings"

// BuildKubectlArgs appends selected context/namespace to a kubectl argument list.
// Use this for all kubectl invocations from this CLI (except context/namespace commands).
func BuildKubectlArgs(args []string) ([]string, error) {
    return BuildKubectlArgsWithOptions(args, "", false)
}

// BuildKubectlArgsWithOptions appends selected context and namespace with overrides.
func BuildKubectlArgsWithOptions(args []string, namespaceOverride string, allNamespaces bool) ([]string, error) {
    contextName, namespace, err := loadKubeSelection()
    if err != nil {
        return nil, err
    }

    if strings.TrimSpace(contextName) != "" {
        args = append(args, "--context", contextName)
    }
    if allNamespaces {
        args = append(args, "--all-namespaces")
        return args, nil
    }

    if strings.TrimSpace(namespaceOverride) != "" {
        args = append(args, "--namespace", namespaceOverride)
        return args, nil
    }

    if strings.TrimSpace(namespace) != "" {
        args = append(args, "--namespace", namespace)
    }

    return args, nil
}
