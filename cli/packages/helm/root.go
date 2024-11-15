package helm

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
)

var settings = cli.New()
var helmCommand = &cobra.Command{
	Use:   "helm",
	Short: "Helm Management Tool",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		timeNow := time.Now().String()
		format = fmt.Sprintf("%s [debug] %s\n", timeNow, format)
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func Execute(rootCmd *cobra.Command) {
  actionConfig := new(action.Configuration)

  helmCommand.PersistentFlags().StringP("namespace", "n", "", "Set namespace")
  helmCommand.PersistentFlags().BoolP("debug", "", false, "Debug output")
  
  helmCommand.MarkFlagRequired("namespace")
  helmCommand.AddCommand(
    NewPurgeCommand(actionConfig),
  )

	rootCmd.AddCommand(helmCommand)
  rootCmd.PersistentFlags()
}
