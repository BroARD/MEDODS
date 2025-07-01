package logging

import (
	"io"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(logDir string, level logrus.Level) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(path.Join(logDir, "all.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return nil, err
	}

	l := logrus.New()
	l.SetReportCaller(true)
	l.SetLevel(level)

	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return path.Base(f.Function) + "()", path.Base(f.File) + ":" + string(rune(f.Line))
		},
	}

	l.SetOutput(io.MultiWriter(logFile, os.Stdout))

	return &Logger{l}, nil
}