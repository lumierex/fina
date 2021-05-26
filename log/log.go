package log

// Log log interface
type Log interface {
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Trace(format string, v ...interface{})
}
