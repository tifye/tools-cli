package winmower

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

func (lf *LogFormatter) Write(data []byte) (int, error) {
	b := bytes.TrimSuffix(data, []byte("\n"))
	lines := bytes.Split(b, []byte("\n"))

	first := lines[0]
	if first[0] == byte('[') {
		switch {
		case bytes.Contains(first, []byte("ERROR")):
			lf.logger.Error(string(trimWinMowerLogLine(first)))
		case bytes.Contains(first, []byte("WARNING")):
			lf.logger.Warn(string(trimWinMowerLogLine(first)))
		case bytes.Contains(first, []byte("INFO")):
			lf.logger.Info(string(trimWinMowerLogLine(first)))
		case bytes.Contains(first, []byte("DEBUG")):
			lf.logger.Debug(string(trimWinMowerLogLine(first)))
		default:
			lf.logger.Info(string(trimWinMowerLogLine(first)))
		}
	}

	rest := lines[1:]
	for _, line := range rest {
		lf.Write(line)
	}

	return len(data), nil
}

func trimWinMowerLogLine(line []byte) []byte {
	const logLevelPrefixLen = 9 // winmower prefixes log level
	return line[logLevelPrefixLen:]
}
