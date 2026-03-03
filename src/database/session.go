package database

import (
    "errors"
    "fmt"
    "os"
    "os/exec"
    "sort"
    "strconv"
    "strings"

    "ehvgo/src/config"
    "ehvgo/src/ui"

    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

type dbSessionOption struct {
    ID        string
    Endpoint  string
    Port      int
    LocalPort int
    Instance  string
    Profile   string
}

func newSessionCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "session",
        Short: "Manage database sessions",
    }

    cmd.AddCommand(newSessionStartCommand())
    return cmd
}

func newSessionStartCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "start",
        Short: "Start an SSM session to a configured database",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            options, err := loadDatabaseSessions()
            if err != nil {
                return err
            }

            selected, err := promptForDatabaseSession(options)
            if err != nil {
                return handlePromptErr(err)
            }

            printSelection(cmd.OutOrStdout(), "Profile", selected.Profile)
            printSelection(cmd.OutOrStdout(), "Instance", selected.Instance)
            printSelection(cmd.OutOrStdout(), "Database", selected.ID)
            printSelection(cmd.OutOrStdout(), "Endpoint", selected.Endpoint)
            printSelection(cmd.OutOrStdout(), "Target port", strconv.Itoa(selected.Port))
            printSelection(cmd.OutOrStdout(), "Local port", strconv.Itoa(selected.LocalPort))

            return startSSMSession(cmd, selected)
        },
    }

    return cmd
}

func loadDatabaseSessions() ([]dbSessionOption, error) {
    cfg, err := config.Read()
    if err != nil {
        return nil, err
    }
    if len(cfg.Databases) == 0 {
        return nil, errors.New("no databases configured; run 'ehvgo db config add'")
    }

    options := make([]dbSessionOption, 0, len(cfg.Databases))
    for id, entry := range cfg.Databases {
        if strings.TrimSpace(id) == "" {
            continue
        }
        if strings.TrimSpace(entry.InstanceID) == "" || strings.TrimSpace(entry.Endpoint) == "" || entry.Port == 0 {
            continue
        }
        if entry.LocalPort == 0 {
            entry.LocalPort = entry.Port
        }
        options = append(options, dbSessionOption{
            ID:        id,
            Endpoint:  entry.Endpoint,
            Port:      entry.Port,
            LocalPort: entry.LocalPort,
            Instance:  entry.InstanceID,
            Profile:   entry.AwsProfile,
        })
    }

    if len(options) == 0 {
        return nil, errors.New("no valid database entries found in config")
    }

    sort.Slice(options, func(i, j int) bool {
        return options[i].ID < options[j].ID
    })

    return options, nil
}

func promptForDatabaseSession(options []dbSessionOption) (dbSessionOption, error) {
    selectPrompt := promptui.Select{
        Label:  "Select database",
        Items:  options,
        Size:   10,
        Stdout: bellSkipper{},
        Templates: &promptui.SelectTemplates{
            Active:   "  \x1b[4m{{ .ID }} | {{ .Endpoint }}:{{ .Port }} ({{ .Instance }})\x1b[0m",
            Inactive: "  {{ .ID }} | {{ .Endpoint }}:{{ .Port }} ({{ .Instance }})",
            Selected: "  {{ .ID }}",
        },
        HideSelected: true,
    }

    index, _, err := selectPrompt.Run()
    if err != nil {
        return dbSessionOption{}, err
    }

    return options[index], nil
}

func startSSMSession(cmd *cobra.Command, selection dbSessionOption) error {
    if strings.TrimSpace(selection.Profile) == "" {
        return errors.New("aws profile is required for session start")
    }

    parameters := fmt.Sprintf("host=%s,portNumber=%d,localPortNumber=%d", selection.Endpoint, selection.Port, selection.LocalPort)
    args := []string{
        "ssm",
        "start-session",
        "--profile", selection.Profile,
        "--target", selection.Instance,
        "--document-name", "AWS-StartPortForwardingSessionToRemoteHost",
        "--parameters", parameters,
    }

    execCmd := exec.CommandContext(cmd.Context(), "aws", args...)
    execCmd.Stdin = cmd.InOrStdin()

    stopSpinner := ui.StartSpinner(os.Stderr, "Starting SSM session")
    execCmd.Stdout = ui.WrapWriterOnFirstWrite(cmd.OutOrStdout(), stopSpinner)
    execCmd.Stderr = ui.WrapWriterOnFirstWrite(cmd.ErrOrStderr(), stopSpinner)

    err := execCmd.Run()
    stopSpinner()
    if err != nil {
        if errors.Is(err, exec.ErrNotFound) {
            return errors.New("aws CLI not found in PATH")
        }
        return err
    }

    return nil
}
