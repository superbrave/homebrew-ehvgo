package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ehvg",
	Short: "A simple tool to make the EHVG dev life a bit easier",
	Long:  "A longer description of the short description to tell you this makes your dev life at EHVG a bit easier",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
