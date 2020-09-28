package log

// LoggerI defines the different log level operations to be implemented by the logger
type LoggerI interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Logger instance
var Logger LoggerI

// Init initializes the logger instance
func Init() {
	Logger = initZap()
}
