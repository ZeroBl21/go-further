package jsonlog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	LevelFatal slog.Level = 12
)

type Logger struct {
	*slog.Logger
}

func New(handler io.Writer, minLevel slog.Level) *Logger {
	logger := slog.New(slog.NewJSONHandler(handler, &slog.HandlerOptions{
		Level:       minLevel,
		ReplaceAttr: replaceAttr,
	}))

	return &Logger{
		Logger: logger,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, slog.Group("properties", args...))
}

func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, slog.Group("properties", args...))
}

func (l *Logger) Error(msg error, args ...any) {
	l.Logger.Error(msg.Error(), slog.Group("properties", args...), formatStackTrace(debug.Stack()))
}

func (l *Logger) Fatal(err error, args ...any) {
	l.Log(
		context.Background(),
		LevelFatal,
		err.Error(),
		slog.Group("properties", args...),
		formatStackTrace(debug.Stack()),
	)
	os.Exit(1) // For entries at the FATAL level, we also terminate the application.
}

func (l *Logger) Write(p []byte) (n int, err error) {
	message := string(p)
	l.Logger.Error(message)
	return len(p), nil
}

func formatStackTrace(stack []byte) slog.Attr {
	slice := strings.Split(string(stack), "\n")

	var trace []any
	for i, v := range slice {
		trace = append(trace, slog.String(strconv.Itoa(i), v))
	}

	return slog.Group("trace", trace[:len(trace)-1]...)
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key != slog.TimeKey {
		return a
	}

	t := a.Value.Time()
	a.Value = slog.StringValue(t.Format(time.RFC3339))

	return a
}
