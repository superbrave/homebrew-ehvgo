package database

import (
    "errors"
    "fmt"
    "sort"
    "strings"
    "text/tabwriter"

    "ehvgo/src/config"

    "github.com/spf13/cobra"
)

func newConfigListCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List configured databases",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg, err := config.Read()
            if err != nil {
                return err
            }
            if len(cfg.Databases) == 0 {
                return errors.New("no databases configured")
            }

            entries := make([]dbSessionOption, 0, len(cfg.Databases))
            for id, entry := range cfg.Databases {
                if strings.TrimSpace(id) == "" {
                    continue
                }
                entries = append(entries, dbSessionOption{
                    ID:        id,
                    Endpoint:  entry.Endpoint,
                    Port:      entry.Port,
                    LocalPort: entry.LocalPort,
                    Instance:  entry.InstanceID,
                    Profile:   entry.AwsProfile,
                })
            }

            sort.Slice(entries, func(i, j int) bool {
                return entries[i].ID < entries[j].ID
            })

            w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
            fmt.Fprintln(w, "ID\tENDPOINT\tPORT\tLOCAL_PORT\tINSTANCE_ID\tPROFILE")
            for _, entry := range entries {
                fmt.Fprintf(
                    w,
                    "%s\t%s\t%d\t%d\t%s\t%s\n",
                    entry.ID,
                    entry.Endpoint,
                    entry.Port,
                    entry.LocalPort,
                    entry.Instance,
                    entry.Profile,
                )
            }
            return w.Flush()
        },
    }

    return cmd
}
