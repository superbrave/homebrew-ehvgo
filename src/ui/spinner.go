package ui

import (
    "fmt"
    "io"
    "os"
    "strings"
    "sync"
    "time"

    "golang.org/x/term"
)

// StartSpinner starts a simple terminal spinner and returns a stop function.
// It is a no-op when stderr is not a terminal.
func StartSpinner(out *os.File, message string) func() {
    if !spinnerEnabled() {
        return func() {}
    }
    if out == nil || !term.IsTerminal(int(out.Fd())) {
        return func() {}
    }

    stop := make(chan struct{})
    done := make(chan struct{})
    frames := []string{"|", "/", "-", "\\"}

    go func() {
        ticker := time.NewTicker(120 * time.Millisecond)
        defer ticker.Stop()
        defer close(done)

        i := 0
        for {
            select {
            case <-stop:
                clear := "\r" + strings.Repeat(" ", len(message)+2) + "\r"
                fmt.Fprint(out, clear)
                return
            case <-ticker.C:
                frame := frames[i%len(frames)]
                fmt.Fprintf(out, "\r%s %s", message, frame)
                i++
            }
        }
    }()

    return func() {
        close(stop)
        <-done
    }
}

// RunWithSpinner runs a function while showing a spinner.
func RunWithSpinner(out *os.File, message string, fn func() error) error {
    stop := StartSpinner(out, message)
    defer stop()
    return fn()
}

// WrapWriterOnFirstWrite stops a spinner when the wrapped writer receives output.
func WrapWriterOnFirstWrite(w io.Writer, stop func()) io.Writer {
    if w == nil {
        return w
    }
    return &spinnerWriter{inner: w, stop: stop}
}

func spinnerEnabled() bool {
    value := strings.TrimSpace(strings.ToLower(os.Getenv("EHVG_SPINNER")))
    return value == "1" || value == "true" || value == "yes" || value == "on"
}

type spinnerWriter struct {
    inner io.Writer
    stop  func()
    once  sync.Once
}

func (sw *spinnerWriter) Write(p []byte) (int, error) {
    if sw.stop != nil {
        sw.once.Do(sw.stop)
    }
    return sw.inner.Write(p)
}
