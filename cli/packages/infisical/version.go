package infisical

import (
	"ehvg/packages/util"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of the EHVGO command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("EHVGo Version " + util.EHVGO_VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
