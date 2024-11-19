package util

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

var  (
  EhvgoVersion string

  Red = color.New(color.FgHiRed)
  Yellow = color.New(color.FgHiYellow)
  Green = color.New(color.FgHiGreen)
  White = color.New(color.FgHiWhite)
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

func GetConfigDir() string {
	h, _ := os.UserHomeDir()
	dir := strings.Join([]string{h, ".ehvg"}, string(os.PathSeparator))

	return dir
}