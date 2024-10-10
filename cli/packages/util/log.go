package util

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintMessageAndExit(messages ...string) {
	if len(messages) > 0 {
		for _, message := range messages {
			fmt.Fprintln(os.Stderr, message)
		}
	}

	os.Exit(1)
}

func PrintError(e string) {
	color.New(color.FgRed).Fprintf(os.Stderr, "error:%v\n", e)
}
