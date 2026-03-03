package database

import (
    "errors"
    "fmt"
    "strings"

    "ehvgo/src/config"

    "github.com/spf13/cobra"
)

func newConfigEditCommand() *cobra.Command {
    var endpoint string
    var instanceID string
    var profile string
    var port int
    var localPort int

    cmd := &cobra.Command{
        Use:   "edit <database-id>",
        Short: "Edit a configured database",
        Args: func(cmd *cobra.Command, args []string) error {
            if len(args) == 0 {
                return errors.New("database ID is required")
            }
            if len(args) > 1 {
                return errors.New("only one database ID is allowed")
            }
            return nil
        },
        RunE: func(cmd *cobra.Command, args []string) error {
            id := strings.TrimSpace(args[0])
            if id == "" {
                return errors.New("database ID is required")
            }

            changed := cmd.Flags().Changed("endpoint") ||
                cmd.Flags().Changed("port") ||
                cmd.Flags().Changed("local_port") ||
                cmd.Flags().Changed("instance_id") ||
                cmd.Flags().Changed("profile")
            if !changed {
                return errors.New("provide at least one of --endpoint, --port, --local_port, --instance_id, --profile")
            }

            cfg, err := config.Read()
            if err != nil {
                return err
            }

            entry, ok := cfg.Databases[id]
            if !ok {
                return fmt.Errorf("database %q not found", id)
            }

            if cmd.Flags().Changed("endpoint") {
                value := strings.TrimSpace(endpoint)
                if value == "" {
                    return errors.New("endpoint cannot be empty")
                }
                entry.Endpoint = value
            }
            if cmd.Flags().Changed("port") {
                if port <= 0 || port > 65535 {
                    return errors.New("port must be between 1 and 65535")
                }
                entry.Port = port
            }
            if cmd.Flags().Changed("local_port") {
                if localPort <= 0 || localPort > 65535 {
                    return errors.New("local_port must be between 1 and 65535")
                }
                if isLocalPortTaken(cfg, id, localPort) {
                    return fmt.Errorf("local_port %d is already used by another database", localPort)
                }
                entry.LocalPort = localPort
            }
            if cmd.Flags().Changed("instance_id") {
                value := strings.TrimSpace(instanceID)
                if value == "" {
                    return errors.New("instance_id cannot be empty")
                }
                entry.InstanceID = value
            }
            if cmd.Flags().Changed("profile") {
                value := strings.TrimSpace(profile)
                if value == "" {
                    return errors.New("profile cannot be empty")
                }
                entry.AwsProfile = value
            }

            cfg.Databases[id] = entry
            if err := config.Write(cfg); err != nil {
                return err
            }

            printSelection(cmd.OutOrStdout(), "Updated", id)
            return nil
        },
    }

    cmd.Flags().StringVar(&endpoint, "endpoint", "", "Database endpoint")
    cmd.Flags().IntVar(&port, "port", 0, "Database port")
    cmd.Flags().IntVar(&localPort, "local_port", 0, "Local port")
    cmd.Flags().StringVar(&instanceID, "instance_id", "", "Jump host instance ID")
    cmd.Flags().StringVar(&profile, "profile", "", "AWS profile")

    return cmd
}

func isLocalPortTaken(cfg config.AppConfig, currentID string, localPort int) bool {
    for id, entry := range cfg.Databases {
        if id == currentID {
            continue
        }
        if entry.LocalPort == localPort {
            return true
        }
    }
    return false
}
