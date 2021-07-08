package easyconfmgr

import (
	"fmt"
	"io"
	"log"
	"os"
)

var DefaultLogger = NewSampleLogger(os.Stderr)
var DiscardLogger = &NopLogger{}

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

type NopLogger struct{}

func (logger *NopLogger) Debug(args ...interface{}) {}

func (logger *NopLogger) Debugf(template string, args ...interface{}) {}

func (logger *NopLogger) Info(args ...interface{}) {}

func (logger *NopLogger) Infof(template string, args ...interface{}) {}

func (logger *NopLogger) Warn(args ...interface{}) {}

func (logger *NopLogger) Warnf(template string, args ...interface{}) {}

func (logger *NopLogger) Error(args ...interface{}) {}

func (logger *NopLogger) Errorf(template string, args ...interface{}) {}

type SampleLogger struct {
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

func (logger *SampleLogger) Debug(args ...interface{}) {
	logger.log(logger.DebugLogger, "DEBUG", fmt.Sprint(args...))
}

func (logger *SampleLogger) Debugf(template string, args ...interface{}) {
	logger.log(logger.DebugLogger, "DEBUG", fmt.Sprintf(template, args...))
}

func (logger *SampleLogger) Info(args ...interface{}) {
	logger.log(logger.InfoLogger, "INFO", fmt.Sprint(args...))
}

func (logger *SampleLogger) Infof(template string, args ...interface{}) {
	logger.log(logger.InfoLogger, "INFO", fmt.Sprintf(template, args...))
}

func (logger *SampleLogger) Warn(args ...interface{}) {
	logger.log(logger.WarnLogger, "WARN", fmt.Sprint(args...))
}

func (logger *SampleLogger) Warnf(template string, args ...interface{}) {
	logger.log(logger.WarnLogger, "WARN", fmt.Sprintf(template, args...))
}

func (logger *SampleLogger) Error(args ...interface{}) {
	logger.log(logger.ErrorLogger, "ERROR", fmt.Sprint(args...))
}

func (logger *SampleLogger) Errorf(template string, args ...interface{}) {
	logger.log(logger.ErrorLogger, "ERROR", fmt.Sprintf(template, args...))
}

func (logger *SampleLogger) log(loggr *log.Logger, level string, msg string) {
	if logger == nil {
		return
	}
	loggr.Output(3, fmt.Sprintf(`%s "%s"`, level, msg))
}

func NewSampleLogger(out io.Writer) Logger {
	l := log.New(out, "[easyconfmgr] ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	stdLogger := &SampleLogger{
		InfoLogger:  l,
		DebugLogger: l,
		WarnLogger:  l,
		ErrorLogger: l,
	}
	return stdLogger
}
