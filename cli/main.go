package main

import (
	"ehvg/packages/cmd"
	"ehvg/packages/helm"
	"ehvg/packages/infisical"
	"ehvg/packages/m365"
	"ehvg/packages/util"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "ehvgo",
	Short:            "A simple tool to make the EHVG dev life a bit easier",
	Long:             "A longer description of the short description to tell you this makes your dev life at EHVG a bit easier",
	TraverseChildren: true,
}

func main() {
	defer os.Exit(0)

	cmd.Execute(rootCmd)
	infisical.Execute(rootCmd)
	m365.Execute(rootCmd)
	helm.Execute(rootCmd)
	err := rootCmd.Execute()
	util.HandleError(err, true)
}
