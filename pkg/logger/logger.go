//nolint:gomnd
package logger

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"log/slog"
	"os"
	"time"
)

type Options struct {
	LogFilePath        string
	BufferSize         int
	BufferFlushTimeout time.Duration
	DebugMode          bool
}

func getDefaultConfig() *Options {
	return &Options{
		LogFilePath:        "./logs.log",
		BufferSize:         1024 * 200,
		BufferFlushTimeout: time.Second * 3,
	}
}

func getLogLevel(isDebug bool) slog.Level {
	if isDebug {
		return slog.LevelDebug
	} else {
		return slog.LevelInfo
	}
}

func New(opts *Options) (*os.File, *slog.Logger, error) {
	if opts == nil {
		opts = getDefaultConfig()
	}

	logLevel := getLogLevel(opts.DebugMode)

	f, err := os.OpenFile(opts.LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	writer := bufio.NewWriterSize(f, opts.BufferSize)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	multiWriter := io.MultiWriter(os.Stdout, writer)

	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}))

	slog.SetDefault(logger)

	go func() {
		for {
			time.Sleep(opts.BufferFlushTimeout)
			err := writer.Flush()
			if err != nil && !errors.Is(err, io.ErrShortWrite) {
				logger.Error(err.Error())
			}
		}
	}()

	return f, logger, nil
}
