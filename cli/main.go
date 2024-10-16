package main

import (
	"ehvg/packages/cmd"
	"os"
)

func main() {
	defer os.Exit(0)

	cmd.Execute()
}
