package testutils

import (
	"io"
	"log/slog"
)

func NewDummyLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
