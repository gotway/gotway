package log

import "github.com/sirupsen/logrus"

type Fields map[string]interface{}

type logrusWrapper struct {
	entry *logrus.Entry
}

func (lw *logrusWrapper) Trace(args ...interface{}) {
	lw.entry.Trace(args...)
}
func (lw *logrusWrapper) Debug(args ...interface{}) {
	lw.entry.Debug(args...)
}
func (lw *logrusWrapper) Print(args ...interface{}) {
	lw.entry.Print(args...)
}
func (lw *logrusWrapper) Info(args ...interface{}) {
	lw.entry.Info(args...)
}
func (lw *logrusWrapper) Warn(args ...interface{}) {
	lw.entry.Warn(args...)
}
func (lw *logrusWrapper) Warning(args ...interface{}) {
	lw.entry.Warning(args...)
}
func (lw *logrusWrapper) Error(args ...interface{}) {
	lw.entry.Error(args...)
}
func (lw *logrusWrapper) Fatal(args ...interface{}) {
	lw.entry.Fatal(args...)
}
func (lw *logrusWrapper) Panic(args ...interface{}) {
	lw.entry.Panic(args...)
}

func (lw *logrusWrapper) Tracef(format string, args ...interface{}) {
	lw.entry.Tracef(format, args...)
}
func (lw *logrusWrapper) Debugf(format string, args ...interface{}) {
	lw.entry.Debugf(format, args...)
}
func (lw *logrusWrapper) Printf(format string, args ...interface{}) {
	lw.entry.Printf(format, args...)
}
func (lw *logrusWrapper) Infof(format string, args ...interface{}) {
	lw.entry.Infof(format, args...)
}
func (lw *logrusWrapper) Warnf(format string, args ...interface{}) {
	lw.entry.Warnf(format, args...)
}
func (lw *logrusWrapper) Warningf(format string, args ...interface{}) {
	lw.entry.Warningf(format, args...)
}
func (lw *logrusWrapper) Errorf(format string, args ...interface{}) {
	lw.entry.Errorf(format, args...)
}
func (lw *logrusWrapper) Fatalf(format string, args ...interface{}) {
	lw.entry.Fatalf(format, args...)
}
func (lw *logrusWrapper) Panicf(format string, args ...interface{}) {
	lw.entry.Panicf(format, args...)
}

func (lw *logrusWrapper) Traceln(args ...interface{}) {
	lw.entry.Traceln(args...)
}
func (lw *logrusWrapper) Debugln(args ...interface{}) {
	lw.entry.Debugln(args...)
}
func (lw *logrusWrapper) Println(args ...interface{}) {
	lw.entry.Println(args...)
}
func (lw *logrusWrapper) Infoln(args ...interface{}) {
	lw.entry.Infoln(args...)
}
func (lw *logrusWrapper) Warnln(args ...interface{}) {
	lw.entry.Warnln(args...)
}
func (lw *logrusWrapper) Warningln(args ...interface{}) {
	lw.entry.Warningln(args...)
}
func (lw *logrusWrapper) Errorln(args ...interface{}) {
	lw.entry.Errorln(args...)
}
func (lw *logrusWrapper) Fatalln(args ...interface{}) {
	lw.entry.Fatalln(args...)
}
func (lw *logrusWrapper) Panicln(args ...interface{}) {
	lw.entry.Panicln(args...)
}

func (lw *logrusWrapper) WithField(key string, value interface{}) Logger {
	return newLogrusWrapper(lw.entry.WithField(key, value))
}

func (lw *logrusWrapper) WithFields(fields Fields) Logger {
	return newLogrusWrapper(lw.entry.WithFields(logrus.Fields(fields)))
}

func newLogrusWrapper(entry *logrus.Entry) Logger {
	return &logrusWrapper{
		entry: entry,
	}
}
