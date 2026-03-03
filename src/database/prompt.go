package database

import (
    "fmt"
    "io"
    "os"
    "errors"

    "github.com/manifoldco/promptui"
    "golang.org/x/term"
)

type bellSkipper struct{}

func (bellSkipper) Write(p []byte) (int, error) {
    filtered := make([]byte, 0, len(p))
    for _, b := range p {
        if b != '\a' {
            filtered = append(filtered, b)
        }
    }
    return os.Stdout.Write(filtered)
}

func (bellSkipper) Close() error {
    return nil
}

func selectTemplates() *promptui.SelectTemplates {
    return &promptui.SelectTemplates{
        Active:   "  {{ . }}",
        Inactive: "  {{ . }}",
        Selected: "",
    }
}

func printSelection(out io.Writer, label, value string) {
    if out == nil {
        return
    }
    if term.IsTerminal(int(os.Stdout.Fd())) {
        const (
            colorGreen = "\x1b[32m"
            colorGray  = "\x1b[90m"
            colorReset = "\x1b[0m"
        )
        fmt.Fprintf(out, "%s✔%s %s%s: %s%s\n", colorGreen, colorReset, colorGray, label, value, colorReset)
        return
    }

    fmt.Fprintf(out, "✔ %s: %s\n", label, value)
}

func handlePromptErr(err error) error {
    if err == nil {
        return nil
    }
    if errors.Is(err, promptui.ErrAbort) {
        return nil
    }
    if errors.Is(err, promptui.ErrInterrupt) {
        return nil
    }
    if errors.Is(err, promptui.ErrEOF) {
        return nil
    }
    return err
}
