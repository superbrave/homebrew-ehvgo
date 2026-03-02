package kubernetes

import (
    "errors"
    "os"
    "path/filepath"
    "strings"

    kubeclient "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func loadKubeSelectionRequired() (string, string, error) {
    contextName, namespace, err := loadKubeSelection()
    if err != nil {
        return "", "", err
    }
    if strings.TrimSpace(contextName) == "" {
        return "", "", errors.New("no context set; run 'ehvg k8s set-context'")
    }
    if strings.TrimSpace(namespace) == "" {
        return "", "", errors.New("no namespace set; run 'ehvg k8s set-namespace'")
    }

    return contextName, namespace, nil
}

func newClientSet(contextName string) (*kubeclient.Clientset, error) {
    restConfig, err := newRestConfig(contextName)
    if err != nil {
        return nil, err
    }

    return kubeclient.NewForConfig(restConfig)
}

func newRestConfig(contextName string) (*rest.Config, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    kubeconfig := filepath.Join(home, ".kube", "config")
    overrides := &clientcmd.ConfigOverrides{CurrentContext: contextName}
    loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
    clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)

    return clientConfig.ClientConfig()
}
