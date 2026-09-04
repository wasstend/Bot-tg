package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	logsPath         = "logs"
	basicLogFilePerm = 0644
	basicDirPerm     = 0755
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func New() *Logger {
	logger, file := setupLogger()
	return &Logger{
		Logger: logger,
		file:   file,
	}
}

func setupLogger() (*slog.Logger, *os.File) {
	currentTime := time.Now().Format("2006-01-02_15:04:05")

	filePath := filepath.Join(logsPath, fmt.Sprintf("%s.log", currentTime))
	if err := os.MkdirAll(filepath.Dir(filePath), basicDirPerm); err != nil {
		panic("failed to create logs directory")
	}

	file, err := os.OpenFile(
		filePath,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		basicLogFilePerm,
	)
	if err != nil {
		panic("failed to open logs file: " + err.Error())
	}

	writer := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	return slog.New(handler), file
}

func (l *Logger) Close() {
	err := l.file.Close()
	if err != nil {
		panic("failed to close logs file: " + err.Error())
	}
}

func (l *Logger) Fatal(msg string, v ...any) {
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Fatal] to capture the real caller

	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	r.Add(v...)

	_ = l.Handler().Handle(context.Background(), r)

	os.Exit(1)
}
