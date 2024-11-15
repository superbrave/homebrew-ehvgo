package helm

import (
	"time"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	h "helm.sh/helm/v3/pkg/time"
)

type EHVGRelease struct {
  *release.Release
}

func filterReleases(cfg *action.Configuration, cmd *cobra.Command, args []string) ([]EHVGRelease, error) {
  list := action.NewList(cfg)
  list.All = true
  list.Short = true
  list.Filter, _ = cmd.Flags().GetString("filter")

  res, err := list.Run()
  if err != nil {
    return nil, err
  }

  releases := []EHVGRelease{}

  if len(res) > 0 {
    for _, r := range res {
      d := EHVGRelease{
        Release: r,
      }

      releases = append(releases, d)
    }
  }

  return releases, nil
}

func (r EHVGRelease) isStale(maxAge int) bool {
  hCutoff := h.Time{
    Time: time.Now().AddDate(0, 0, -(maxAge)),
  }

  return r.Info.LastDeployed.Before(hCutoff)
}