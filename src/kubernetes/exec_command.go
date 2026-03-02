package kubernetes

import (
    "context"
    "errors"
    "io"
    "os"
    "os/signal"
    "sort"
    "strings"
    "sync"
    "syscall"
    "time"

	"ehvgo/src/ui"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func newExecCommand() *cobra.Command {
	var usePod bool
	var useDeployment bool
	var shellCommand string

	cmd := &cobra.Command{
		Use:   "exec [container]",
		Short: "Execute a shell in a container",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if usePod == useDeployment {
				return errors.New("exactly one of --pod or --deployment is required")
			}

			container := ""
			if len(args) == 1 {
				container = strings.TrimSpace(args[0])
			}

			contextName, namespace, err := loadKubeSelectionRequired()
			if err != nil {
				return err
			}

			printContextAndNamespace(cmd.OutOrStdout(), contextName, namespace)

			clientSet, err := newClientSet(contextName)
			if err != nil {
				return err
			}
			restConfig, err := newRestConfig(contextName)
			if err != nil {
				return err
			}

			var target string
			if usePod {
				var pods []string
				err = ui.RunWithSpinner(os.Stderr, "Fetching pods", func() error {
					var listErr error
					pods, listErr = listPods(clientSet, namespace)
					return listErr
				})
				if err != nil {
					return err
				}

				selected, err := selectItem("Select pod", pods)
				if err != nil {
					return err
				}
				if strings.TrimSpace(selected) == "" {
					return errors.New("pod selection is required")
				}
				target = selected
				printSelection(cmd.OutOrStdout(), "Pod", target)

                container, err = ensureContainerSelected(cmd.OutOrStdout(), container, "Fetching containers", func() ([]string, error) {
                    return listPodContainers(clientSet, namespace, target)
                })
                if err != nil {
                    return err
                }
            }

			if useDeployment {
				var deployments []string
				err = ui.RunWithSpinner(os.Stderr, "Fetching deployments", func() error {
					var listErr error
					deployments, listErr = listDeployments(clientSet, namespace)
					return listErr
				})
				if err != nil {
					return err
				}

				selected, err := selectItem("Select deployment", deployments)
				if err != nil {
					return err
				}
				if strings.TrimSpace(selected) == "" {
					return errors.New("deployment selection is required")
				}

				var podName string
				err = ui.RunWithSpinner(os.Stderr, "Selecting pod", func() error {
					var selectErr error
					podName, selectErr = selectDeploymentPod(clientSet, namespace, selected)
					return selectErr
				})
				if err != nil {
					return err
				}
				if strings.TrimSpace(podName) == "" {
					return errors.New("no pod available for deployment")
				}

				target = podName
				printSelection(cmd.OutOrStdout(), "Deployment", selected)
				printSelection(cmd.OutOrStdout(), "Pod", target)

                container, err = ensureContainerSelected(cmd.OutOrStdout(), container, "Fetching containers", func() ([]string, error) {
                    return listDeploymentContainers(clientSet, namespace, selected)
                })
                if err != nil {
                    return err
                }
            }

			command := strings.TrimSpace(shellCommand)
			if command == "" {
				command = "bash"
			}

			stopSpinner := ui.StartSpinner(os.Stderr, "Starting shell")
			stopSpinner()

			err = execIntoPod(restConfig, namespace, target, container, strings.Fields(command))
			if err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&usePod, "pod", "p", false, "Select a pod")
	cmd.Flags().BoolVarP(&useDeployment, "deployment", "d", false, "Select a deployment")
	cmd.Flags().StringVar(&shellCommand, "command", "bash", "Command to run in the container")

	ui.AddHelpCommand(cmd)
	return cmd
}

func selectItem(label string, items []string) (string, error) {
	if len(items) == 0 {
		return "", errors.New("no items available")
	}

	selectPrompt := promptui.Select{
		Label:        label,
		Items:        items,
		Size:         10,
		Stdout:       bellSkipper{},
		Templates:    selectTemplates(),
		HideSelected: true,
	}

	_, result, err := selectPrompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrAbort) {
			return "", errors.New("selection cancelled")
		}
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func listPods(clientSet *kubeclient.Clientset, namespace string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		return nil, errors.New("no pods found in namespace")
	}

	pods := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		name := strings.TrimSpace(item.Name)
		if name != "" {
			pods = append(pods, name)
		}
	}
	sort.Strings(pods)

	return pods, nil
}

func listDeployments(clientSet *kubeclient.Clientset, namespace string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	list, err := clientSet.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if len(list.Items) == 0 {
		return nil, errors.New("no deployments found in namespace")
	}

	deployments := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		name := strings.TrimSpace(item.Name)
		if name != "" {
			deployments = append(deployments, name)
		}
	}
	sort.Strings(deployments)

	return deployments, nil
}

func execIntoPod(restConfig *rest.Config, namespace, podName, container string, command []string) error {
	if restConfig == nil {
		return errors.New("kubernetes config is not available")
	}

	cfg := rest.CopyConfig(restConfig)
	cfg.GroupVersion = &corev1.SchemeGroupVersion
	cfg.APIPath = "/api"
	cfg.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	restClient, err := rest.RESTClientFor(cfg)
	if err != nil {
		return err
	}

	req := restClient.Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	tty := stdin != nil && term.IsTerminal(int(stdin.Fd()))
	if tty {
		oldState, err := term.MakeRaw(int(stdin.Fd()))
		if err != nil {
			return err
		}
        defer func() {
            _ = term.Restore(int(stdin.Fd()), oldState)
        }()
	}

	sizeQueue := newTerminalSizeQueue()
	sizeQueue.start(stdin)
	execOptions := &corev1.PodExecOptions{
		Container: container,
		Command:   command,
		Stdin:     true,
		Stdout:    stdout != nil,
		Stderr:    stderr != nil && !tty,
		TTY:       tty,
	}

	req.VersionedParams(execOptions, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		return err
	}

	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             stdin,
		Stdout:            stdout,
		Stderr:            stderr,
		Tty:               tty,
		TerminalSizeQueue: sizeQueue,
	})
}

func ensureContainerSelected(out io.Writer, current, spinnerMessage string, listFn func() ([]string, error)) (string, error) {
    if strings.TrimSpace(current) != "" {
        return current, nil
    }

    var containers []string
    err := ui.RunWithSpinner(os.Stderr, spinnerMessage, func() error {
        var listErr error
        containers, listErr = listFn()
        return listErr
    })
    if err != nil {
        return "", err
    }

    selected, err := selectItem("Select container", containers)
    if err != nil {
        return "", err
    }
    if strings.TrimSpace(selected) == "" {
        return "", errors.New("container selection is required")
    }

    printSelection(out, "Container", selected)
    return selected, nil
}

type terminalSizeQueue struct {
	once   sync.Once
	sizes  chan remotecommand.TerminalSize
	closed chan struct{}
}

func newTerminalSizeQueue() *terminalSizeQueue {
	return &terminalSizeQueue{
		sizes:  make(chan remotecommand.TerminalSize, 1),
		closed: make(chan struct{}),
	}
}

func (q *terminalSizeQueue) Next() *remotecommand.TerminalSize {
	select {
	case size := <-q.sizes:
		return &size
	case <-q.closed:
		return nil
	}
}

func (q *terminalSizeQueue) start(stdin *os.File) {
	q.once.Do(func() {
		if stdin == nil || !term.IsTerminal(int(stdin.Fd())) {
			close(q.closed)
			return
		}

		q.pushSize(stdin)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGWINCH)

		go func() {
			defer signal.Stop(sigs)
			defer close(q.closed)
			for range sigs {
				q.pushSize(stdin)
			}
		}()
	})
}

func (q *terminalSizeQueue) pushSize(stdin *os.File) {
	width, height, err := term.GetSize(int(stdin.Fd()))
	if err != nil {
		return
	}
	size := remotecommand.TerminalSize{Width: uint16(width), Height: uint16(height)}
	select {
	case q.sizes <- size:
	default:
		// drop if channel is full
	}
}

func selectDeploymentPod(clientSet *kubeclient.Clientset, namespace, deploymentName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	if deployment.Spec.Selector == nil {
		return "", errors.New("deployment has no selector")
	}

	selector := metav1.FormatLabelSelector(deployment.Spec.Selector)
	if strings.TrimSpace(selector) == "" {
		return "", errors.New("deployment selector is empty")
	}

	pods, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return "", err
	}

	if len(pods.Items) == 0 {
		return "", errors.New("no pods found for deployment")
	}

	sorted := pods.Items
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Name < sorted[j].Name
	})

	for _, pod := range sorted {
		if pod.Status.Phase == corev1.PodRunning && isPodReady(pod) {
			return pod.Name, nil
		}
	}
	for _, pod := range sorted {
		if pod.Status.Phase == corev1.PodRunning {
			return pod.Name, nil
		}
	}
	return sorted[0].Name, nil
}

func isPodReady(pod corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func listPodContainers(clientSet *kubeclient.Clientset, namespace, podName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pod, err := clientSet.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(pod.Spec.Containers) == 0 {
		return nil, errors.New("no containers found in pod")
	}

	containers := make([]string, 0, len(pod.Spec.Containers))
	for _, container := range pod.Spec.Containers {
		name := strings.TrimSpace(container.Name)
		if name != "" {
			containers = append(containers, name)
		}
	}
	sort.Strings(containers)

	return containers, nil
}

func listDeploymentContainers(clientSet *kubeclient.Clientset, namespace, deploymentName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deployment, err := clientSet.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if len(deployment.Spec.Template.Spec.Containers) == 0 {
		return nil, errors.New("no containers found in deployment")
	}

	containers := make([]string, 0, len(deployment.Spec.Template.Spec.Containers))
	for _, container := range deployment.Spec.Template.Spec.Containers {
		name := strings.TrimSpace(container.Name)
		if name != "" {
			containers = append(containers, name)
		}
	}
	sort.Strings(containers)

	return containers, nil
}
