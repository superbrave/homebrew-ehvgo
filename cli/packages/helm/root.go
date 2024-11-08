package helm

import "github.com/spf13/cobra"

var helmCommand = &cobra.Command{
	Use:   "helm",
	Short: "Helm Management Tool",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute(rootCmd *cobra.Command) {
	rootCmd.AddCommand(helmCommand)
}
