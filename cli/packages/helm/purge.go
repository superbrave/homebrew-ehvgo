package helm

import (
	"ehvg/packages/util"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	hAction "helm.sh/helm/v3/pkg/action"
	hCli "helm.sh/helm/v3/pkg/cli"
	hTime "helm.sh/helm/v3/pkg/time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var actionConfig = &hAction.Configuration{}

var purgeCommand = &cobra.Command{
	Use:   "purge",
	Short: "Purge stale deployments",
	Run:   purgeDeployments,
}

func purgeDeployments(cmd *cobra.Command, args []string) {
	namespace, _ := cmd.Flags().GetString("namespace")
	maxAge, _ := cmd.Flags().GetInt("max-age")
	filter, _ := cmd.Flags().GetString("filter")
	showList, _ := cmd.Flags().GetBool("list")
	ack, _ := cmd.Flags().GetBool("yes")

	settings := hCli.New()

	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		fmt.Println(err)
	}

	deployments := hAction.NewList(actionConfig)
	deployments.All = true
	deployments.Filter = filter
	deployments.Short = true

	list, err := deployments.Run()
	if err != nil {
		color.New(color.FgHiRed).Printf("Cannot list deployments: %v", err)

		return
	}

	cutoff := time.Now().AddDate(0, 0, -(maxAge))
	hCutoff := hTime.Time{
		Time: time.Now().AddDate(0, 0, -(maxAge)),
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Revision", "Name", "Installed at", "Updated at", "# Days left"})

	rows := []table.Row{}

	if len(list) <= 0 {
		util.White.Printf("No releases found with the current filter: %v", filter)

		return
	}

	for _, release := range list {
		isStale := release.Info.LastDeployed.Before(hCutoff)
		daysLeft := math.Round(release.Info.LastDeployed.Time.Sub(cutoff).Hours() / 24)
		daysLeftStr := fmt.Sprintf("%.0f", daysLeft)

		c := util.Green.Sprint(daysLeftStr)

		if daysLeft < 1 {
			c = util.Red.Sprint(daysLeftStr)
		}

		if showList {
			rows = append(rows, table.Row{release.Version, release.Name, release.Info.FirstDeployed.Format(time.DateTime), release.Info.LastDeployed.Format(time.DateTime), c})
		}

		if isStale {
			uninstall := hAction.NewUninstall(actionConfig)
			uninstall.DryRun = !ack
			uninstall.KeepHistory = false
			uninstall.DeletionPropagation = "Foreground"
			uninstall.Wait = true
			uninstall.IgnoreNotFound = false

			uninstallRun, err := uninstall.Run(release.Name)
			if err != nil {
				util.Red.Printf("Unable to uninstall release %v: %v", release.Name, err)
				return
			}

			util.White.Printf(uninstallRun.Info)
		}
	}

	t.AppendRows(rows)

	if showList {
		t.Render()
	}

}

func init() {
	purgeCommand.Flags().StringP("namespace", "n", "", "Namespace to look for deployments")
	purgeCommand.Flags().IntP("max-age", "", 14, "Maximum lifetime of a release, before it gets purged")
	purgeCommand.Flags().StringP("filter", "", "", "Filter for deployment names")
	purgeCommand.Flags().BoolP("list", "l", false, "List all the releases")
	purgeCommand.Flags().BoolP("yes", "y", false, "Explicitly set flag to acknowledge that you want to remove releases")
	purgeCommand.MarkFlagsRequiredTogether("namespace", "max-age", "filter")
	helmCommand.AddCommand(purgeCommand)
}
