package aws

import (
    "bufio"
    "errors"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
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
                    profiles, err := readAWSProfiles()
                    if err != nil {
                        return err
                    }
                    selected, err := promptForProfile(profiles)
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

func readAWSProfiles() ([]string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }

    configPath := filepath.Join(home, ".aws", "config")
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    profiles := make(map[string]struct{})
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
            continue
        }
        section := strings.TrimSpace(line[1 : len(line)-1])
        if section == "default" {
            profiles["default"] = struct{}{}
            continue
        }
        if strings.HasPrefix(section, "profile ") {
            name := strings.TrimSpace(strings.TrimPrefix(section, "profile "))
            if name != "" {
                profiles[name] = struct{}{}
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    if len(profiles) == 0 {
        return nil, errors.New("no profiles found in ~/.aws/config")
    }

    list := make([]string, 0, len(profiles))
    for name := range profiles {
        list = append(list, name)
    }
    sort.Strings(list)

    return list, nil
}

func promptForProfile(profiles []string) (string, error) {
    selectPrompt := promptui.Select{
        Label:  "Select AWS profile",
        Items:  profiles,
        Size:   10,
        Stdout: bellSkipper{},
    }

    _, result, err := selectPrompt.Run()
    if err != nil {
        return "", err
    }

    return result, nil
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
