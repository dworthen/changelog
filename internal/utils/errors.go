package utils

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

func CheckError(err error) {
	if err != nil {
		slog.Error("Error", "err", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func WithStackTrace(err error) error {
	if err == nil {
		return nil
	}
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	return fmt.Errorf("%w\nStack trace: %s", err, buf[:n])
}
