package main

import (
	"ehvg/packages/infisical"
	"os"
)

func main() {
	defer os.Exit(0)

	infisical.Execute()
}
