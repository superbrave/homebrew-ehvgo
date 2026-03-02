package kubernetes

import "os"

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
