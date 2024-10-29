package m365

import "github.com/spf13/cobra"

var m365Command = &cobra.Command{
	Use:   "m365",
	Short: "Microsoft 365 Management Tool",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func Execute(rootCmd *cobra.Command) {
	rootCmd.AddCommand(m365Command)
}
