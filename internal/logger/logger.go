package logger

import (
	"fmt"
	"log/slog"
)

func Info(message string, args ...any) {
	slog.Info(fmt.Sprintf("| %s | %s", message, args))
}

func Error(message string, args ...any) {
	slog.Error(fmt.Sprintf("| %s | %s", message, args))
}
