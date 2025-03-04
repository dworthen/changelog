package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func CheckError(err error) {
	if err != nil {
		slog.Error("Error", "err", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
