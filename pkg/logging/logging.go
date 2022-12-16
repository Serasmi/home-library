package logging

import (
	"fmt"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

var logger Logger

type Logger struct {
	*logrus.Entry
}

func NewLogger() *Logger {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: true,
		FullTimestamp: true,
	}
	l.SetLevel(logrus.DebugLevel)

	logger = Logger{logrus.NewEntry(l)}

	return &logger
}

func GetLogger() *Logger {
	return &logger
}

func (l *Logger) SetLevel(lvl string) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		l.Errorf("unable to parse loglevel: %s", lvl)
		return
	}

	l.Logger.SetLevel(level)
}
