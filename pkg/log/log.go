package log

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Println(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
}

var Log = newLogrusWrapper(logrus.New().WithField("type", "default"))

func NewLogger(fields Fields, env string, level string, transport io.Writer) Logger {
	l := logrus.New()

	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		l.Panicf("cannot parse level %s: %s", level, err)
	}
	l.SetLevel(logrusLevel)

	if env != "local" {
		l.SetFormatter(&logrus.JSONFormatter{})
	}

	l.SetOutput(transport)
	logger := l.WithFields(logrus.Fields(fields))
	return newLogrusWrapper(logger)
}
