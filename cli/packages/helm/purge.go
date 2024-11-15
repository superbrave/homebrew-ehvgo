package helm

import (
	"ehvg/packages/util"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
)

func NewPurgeCommand(cfg *action.Configuration) *cobra.Command {
  cmd := &cobra.Command{
    Use:   "purge",
    Short: "Purge stale deployments",
    Run:   func(cmd *cobra.Command, args []string) {
        debug, _ := cmd.Flags().GetBool("debug")
        ns, _ := cmd.Flags().GetString("namespace")
        if err := cfg.Init(settings.RESTClientGetter(), ns, os.Getenv("HELM_DRIVER"), debug); err != nil {
          util.Red.Printf("error initializing config: %v", err)
          return
        }

        releases, err := filterReleases(cfg, cmd, args)
        if err != nil {
          util.Red.Printf("error while fetching releases: %v", err)
          return
        }

        if len(releases) < 1 {
          util.White.Printf("No releases found with current filters")
          return
        }

        maxAge, _ := cmd.Flags().GetInt("max-age")
        numDeleted := 0

        for _, r := range releases {
          if r.isStale(maxAge) {
            numDeleted++
            if err := r.purgeDeployment(cfg, debug); err != nil {
              util.Red.Printf("error deleting release: %v", err)
              return
            }
          }
        }

        util.White.Printf("Deleted %v stale deployements\n", numDeleted)
    },
  }

  cmd.Flags().StringP("filter", "f", ".*", "Filter your releases by name")
  cmd.Flags().IntP("max-age", "", 14, "Maximum age of a release in days")

  return cmd
}

func (r EHVGRelease) purgeDeployment(cfg *action.Configuration, debug bool) error {
  util.White.Printf("Deleting %v..\n", r.Name)

  rem := action.NewUninstall(cfg)
  if debug {
    util.White.Println(rem)
  }
  rem.DeletionPropagation = "foreground"
  rem.Wait = true
  
  resp, err := rem.Run(r.Name)
  if err != nil {
    return err
  }

  fmt.Println(resp.Info)
 
  return nil

}
