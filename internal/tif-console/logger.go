package tifconsole

import (
	"bytes"

	"github.com/charmbracelet/log"
)

type LogFormatter struct {
	logger *log.Logger
}

func NewLogFormatter(logger *log.Logger) *LogFormatter {
	return &LogFormatter{
		logger: logger,
	}
}

func (l *LogFormatter) Write(data []byte) (int, error) {
	trimmed := bytes.TrimSuffix(data, []byte("\n"))
	l.logger.Print(string(trimmed))
	return len(data), nil
}
