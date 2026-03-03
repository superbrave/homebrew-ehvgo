package aws

import (
    "errors"
    "os"
    "os/exec"
    "strings"

    "ehvgo/src/ui"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

type bellSkipper struct{}

func (bellSkipper) Write(p []byte) (int, error) {
    filtered := make([]byte, 0, len(p))
    for _, b := range p {
        if b != '\a' {
            filtered = append(filtered, b)
        }
    }
    return os.Stdout.Write(filtered)
}

func (bellSkipper) Close() error {
    return nil
}

func newLoginCommand() *cobra.Command {
    var loginProfile string

    cmd := &cobra.Command{
        Use:   "login",
        Short: "Authenticate with AWS SSO",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            profile := strings.TrimSpace(loginProfile)
            if profile == "" {
                envProfile := strings.TrimSpace(os.Getenv("AWS_PROFILE"))
                if envProfile != "" {
                    profile = envProfile
                } else {
                    profiles, err := ReadProfiles()
                    if err != nil {
                        return err
                    }
                    selected, err := PromptForProfile(profiles)
                    if err != nil {
                        return err
                    }
                    profile = selected
                }
            }
            if profile == "" {
                return errors.New("profile is required")
            }

            openBrowser, err := promptYesNo("Open in browser?", false)
            if err != nil {
                return err
            }

            argsList := []string{"sso", "login", "--profile", profile}
            if !openBrowser {
                argsList = append(argsList, "--no-browser")
            }

            execCmd := exec.Command("aws", argsList...)
            execCmd.Stdin = cmd.InOrStdin()

            stopSpinner := ui.StartSpinner(os.Stderr, "Running aws sso login")
            execCmd.Stdout = ui.WrapWriterOnFirstWrite(cmd.OutOrStdout(), stopSpinner)
            execCmd.Stderr = ui.WrapWriterOnFirstWrite(cmd.ErrOrStderr(), stopSpinner)

            err = execCmd.Run()
            stopSpinner()
            if err != nil {
                if errors.Is(err, exec.ErrNotFound) {
                    return errors.New("aws CLI not found in PATH")
                }
                return err
            }

            return nil
        },
    }

    cmd.Flags().StringVar(&loginProfile, "profile", "", "AWS profile name")

    ui.AddHelpCommand(cmd)
    return cmd
}

func promptYesNo(question string, defaultYes bool) (bool, error) {
    defaultValue := "n"
    if defaultYes {
        defaultValue = "y"
    }

    prompt := promptui.Prompt{
        Label:     question,
        IsConfirm: true,
        Default:   defaultValue,
        Stdout:    bellSkipper{},
    }

    result, err := prompt.Run()
    if err != nil {
        if errors.Is(err, promptui.ErrAbort) {
            return false, nil
        }
        return false, err
    }

    normalized := strings.TrimSpace(strings.ToLower(result))
    if normalized == "" {
        normalized = defaultValue
    }

    return normalized == "y" || normalized == "yes", nil
}
