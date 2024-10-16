package util

import (
	"os"
	"runtime"

	"github.com/fatih/color"
)

func PrintError(e error, throwPanic bool) {
	color.New(color.FgRed).Fprintf(os.Stderr, "%v\n", e.Error())

	if throwPanic {
		runtime.Goexit()
	}
}
