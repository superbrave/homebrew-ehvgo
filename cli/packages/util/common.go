package util

import (
	"os"
	"strings"
)

func HandleError(err error, throwPanic bool) {
	if err != nil {
		PrintError(err, throwPanic)
	}
}

func GetCwdForFile(filename string) string {
	wd, err := os.Getwd()
	HandleError(err, true)

	path := []string{wd, strings.Trim(filename, string(os.PathSeparator))}

	return strings.Join(path, string(os.PathSeparator))
}
