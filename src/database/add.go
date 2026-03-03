package database

import (
    "context"
    "errors"
    "fmt"
    "os"
    "strconv"
    "strings"

    awsprofiles "ehvgo/src/aws"
    "ehvgo/src/config"
    "ehvgo/src/ui"

    awssdk "github.com/aws/aws-sdk-go-v2/aws"
    awsconfig "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/ec2"
    ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
    "github.com/aws/aws-sdk-go-v2/service/rds"
    "github.com/manifoldco/promptui"
    "github.com/spf13/cobra"
)

type dbOption struct {
    ID       string
    Endpoint string
    Port     int
    VpcID    string
}

type instanceOption struct {
    ID   string
    Name string
}

func newConfigAddCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "add",
        Short: "Add a database connection profile",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            profiles, err := awsprofiles.ReadProfiles()
            if err != nil {
                return err
            }
            profile, err := awsprofiles.PromptForProfile(profiles)
            if err != nil {
                return handlePromptErr(err)
            }
            if strings.TrimSpace(profile) == "" {
                return errors.New("profile is required")
            }
            printSelection(cmd.OutOrStdout(), "Profile", profile)

            var databases []dbOption
            err = ui.RunWithSpinner(os.Stderr, "Loading databases", func() error {
                var listErr error
                databases, listErr = listDatabases(cmd.Context(), profile)
                return listErr
            })
            if err != nil {
                return err
            }

            selectedDb, err := promptForDatabase(databases)
            if err != nil {
                return handlePromptErr(err)
            }
            printSelection(cmd.OutOrStdout(), "Database", fmt.Sprintf("%s (%s:%d)", selectedDb.ID, selectedDb.Endpoint, selectedDb.Port))

            var instances []instanceOption
            err = ui.RunWithSpinner(os.Stderr, "Loading jump hosts", func() error {
                var listErr error
                instances, listErr = listJumpHosts(cmd.Context(), profile, selectedDb.VpcID)
                return listErr
            })
            if err != nil {
                return err
            }

            selectedInstance, err := promptForInstance(instances)
            if err != nil {
                return handlePromptErr(err)
            }
            label := selectedInstance.ID
            if strings.TrimSpace(selectedInstance.Name) != "" {
                label = fmt.Sprintf("%s (%s)", selectedInstance.Name, selectedInstance.ID)
            }
            printSelection(cmd.OutOrStdout(), "Jump host", label)

            localPort, err := promptForLocalPort(selectedDb.Port)
            if err != nil {
                return handlePromptErr(err)
            }
            printSelection(cmd.OutOrStdout(), "Local port", strconv.Itoa(localPort))

            cfg, err := config.Read()
            if err != nil {
                return err
            }

            cfg.Databases[selectedDb.ID] = config.DatabaseConfig{
                AwsProfile: profile,
                InstanceID: selectedInstance.ID,
                Endpoint:   selectedDb.Endpoint,
                Port:       selectedDb.Port,
                LocalPort:  localPort,
            }

            if err := config.Write(cfg); err != nil {
                return err
            }

            printSelection(cmd.OutOrStdout(), "Saved", selectedDb.ID)
            return nil
        },
    }

    return cmd
}

func promptForDatabase(options []dbOption) (dbOption, error) {
    selectPrompt := promptui.Select{
        Label:        "Select database",
        Items:        options,
        Size:         10,
        Stdout:       bellSkipper{},
        Templates: &promptui.SelectTemplates{
            Active:   "  \x1b[4m{{ .ID }} | {{ .Endpoint }}:{{ .Port }}\x1b[0m",
            Inactive: "  {{ .ID }} | {{ .Endpoint }}:{{ .Port }}",
            Selected: "",
        },
        HideSelected: true,
    }

    index, _, err := selectPrompt.Run()
    if err != nil {
        return dbOption{}, err
    }

    return options[index], nil
}

func promptForInstance(options []instanceOption) (instanceOption, error) {
    selectPrompt := promptui.Select{
        Label:        "Select jump host instance",
        Items:        options,
        Size:         10,
        Stdout:       bellSkipper{},
        Templates: &promptui.SelectTemplates{
            Active:   "  \x1b[4m{{ if .Name }}{{ .Name }} ({{ .ID }}){{ else }}{{ .ID }}{{ end }}\x1b[0m",
            Inactive: "  {{ if .Name }}{{ .Name }} ({{ .ID }}){{ else }}{{ .ID }}{{ end }}",
            Selected: "",
        },
        HideSelected: true,
    }

    index, _, err := selectPrompt.Run()
    if err != nil {
        return instanceOption{}, err
    }

    return options[index], nil
}

func promptForLocalPort(defaultPort int) (int, error) {
    prompt := promptui.Prompt{
        Label:   "Local port",
        Default: strconv.Itoa(defaultPort),
        Stdout:  bellSkipper{},
        Validate: func(input string) error {
            value, err := strconv.Atoi(strings.TrimSpace(input))
            if err != nil {
                return errors.New("enter a valid port number")
            }
            if value < 1 || value > 65535 {
                return errors.New("port must be between 1 and 65535")
            }
            return nil
        },
    }

    result, err := prompt.Run()
    if err != nil {
        return 0, err
    }

    value, err := strconv.Atoi(strings.TrimSpace(result))
    if err != nil {
        return 0, err
    }

    return value, nil
}

func listDatabases(ctx context.Context, profile string) ([]dbOption, error) {
    cfg, err := loadAWSConfig(ctx, profile)
    if err != nil {
        return nil, err
    }

    client := rds.NewFromConfig(cfg)
    output, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{})
    if err != nil {
        return nil, err
    }

    options := make([]dbOption, 0, len(output.DBInstances))
    for _, instance := range output.DBInstances {
        id := strings.TrimSpace(awssdk.ToString(instance.DBInstanceIdentifier))
        if instance.Endpoint == nil {
            continue
        }
        endpoint := strings.TrimSpace(awssdk.ToString(instance.Endpoint.Address))
        port := int(awssdk.ToInt32(instance.Endpoint.Port))
        vpcID := ""
        if instance.DBSubnetGroup != nil {
            vpcID = strings.TrimSpace(awssdk.ToString(instance.DBSubnetGroup.VpcId))
        }
        if id == "" || endpoint == "" || port == 0 || vpcID == "" {
            continue
        }
        options = append(options, dbOption{
            ID:       id,
            Endpoint: endpoint,
            Port:     port,
            VpcID:    vpcID,
        })
    }

    if len(options) == 0 {
        return nil, errors.New("no databases found")
    }

    return options, nil
}

func listJumpHosts(ctx context.Context, profile string, vpcID string) ([]instanceOption, error) {
    cfg, err := loadAWSConfig(ctx, profile)
    if err != nil {
        return nil, err
    }

    client := ec2.NewFromConfig(cfg)
    filters := []ec2types.Filter{
        {
            Name:   awssdk.String("tag:Name"),
            Values: []string{"*jumphost*"},
        },
    }
    if strings.TrimSpace(vpcID) != "" {
        filters = append(filters, ec2types.Filter{
            Name:   awssdk.String("vpc-id"),
            Values: []string{vpcID},
        })
    }

    output, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
        Filters: filters,
    })
    if err != nil {
        return nil, err
    }

    options := make([]instanceOption, 0)
    for _, reservation := range output.Reservations {
        for _, instance := range reservation.Instances {
            id := strings.TrimSpace(awssdk.ToString(instance.InstanceId))
            if id == "" {
                continue
            }
            name := tagValue(instance.Tags, "Name")
            if !isJumpHostName(name) {
                continue
            }
            options = append(options, instanceOption{
                ID:   id,
                Name: name,
            })
        }
    }

    if len(options) == 0 {
        return nil, errors.New("no jump host instances found")
    }

    return options, nil
}

func tagValue(tags []ec2types.Tag, key string) string {
    for _, tag := range tags {
        if strings.EqualFold(strings.TrimSpace(awssdk.ToString(tag.Key)), key) {
            return strings.TrimSpace(awssdk.ToString(tag.Value))
        }
    }
    return ""
}

func isJumpHostName(name string) bool {
    trimmed := strings.TrimSpace(name)
    if trimmed == "" {
        return false
    }
    lowered := strings.ToLower(trimmed)
    return strings.Contains(lowered, "jumphost") || strings.HasSuffix(lowered, "jumphost")
}

func loadAWSConfig(ctx context.Context, profile string) (awssdk.Config, error) {
    cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithSharedConfigProfile(profile))
    if err != nil {
        return awssdk.Config{}, err
    }
    return cfg, nil
}
