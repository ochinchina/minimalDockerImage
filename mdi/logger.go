package main

import (
	"fmt"
	"io"
	"time"
)

type LogLevel int

const (
	Critial LogLevel = iota
	Fatal
	Error
	Warn
	Info
	Debug
)

type Logger struct {
	writer io.Writer
	level  LogLevel
}

func NewLogger(writer io.Writer) *Logger {
	return &Logger{writer: writer, level: Debug}
}

func (lg *Logger) SetOutput(writer io.Writer) {
	if writer != nil {
		lg.writer = writer
	}
}

func (lg *Logger) SetLogLevel(level LogLevel) {
	lg.level = level
}

func (lg *Logger) IsLogEnabled(level LogLevel) bool {
	return level <= lg.level
}

func (lg *Logger) Critical(v ...interface{}) {
	lg.writeLog(Critial, v...)
}

func (lg *Logger) Criticalf(fomat string, v ...interface{}) {
	lg.writeLogf(Critial, fomat, v...)
}

func (lg *Logger) Fatal(v ...interface{}) {
	lg.writeLog(Fatal, v...)
}

func (lg *Logger) Fatalf(format string, v ...interface{}) {
	lg.writeLogf(Fatal, format, v...)
}

func (lg *Logger) Error(v ...interface{}) {
	lg.writeLog(Error, v...)
}

func (lg *Logger) Errorf(format string, v ...interface{}) {
	lg.writeLogf(Error, format, v...)
}

func (lg *Logger) Warn(v ...interface{}) {
	lg.writeLog(Warn, v...)
}

func (lg *Logger) Warnf(format string, v ...interface{}) {
	lg.writeLogf(Warn, format, v...)
}

func (lg *Logger) Info(v ...interface{}) {
	lg.writeLog(Info, v...)
}

func (lg *Logger) Infof(format string, v ...interface{}) {
	lg.writeLogf(Info, format, v...)
}

func (lg *Logger) Debug(v ...interface{}) {
	lg.writeLog(Debug, v...)
}

func (lg *Logger) Debugf(format string, v ...interface{}) {
	lg.writeLogf(Debug, format, v...)
}

func (lg *Logger) writeLog(level LogLevel, v ...interface{}) {
	if lg.IsLogEnabled(level) {
		fmt.Fprintln(lg.writer, v...)
	}
}

func (lg *Logger) writeLogf(level LogLevel, format string, v ...interface{}) {
	if lg.IsLogEnabled(level) {
		lg.writeHeader(level)
		fmt.Fprintf(lg.writer, format, v...)
		fmt.Fprintln(log.writer)
	}
}

func (lg *Logger) writeHeader(level LogLevel) {
	t := time.Now()
	fmt.Fprintf(lg.writer, "%d-%02d-%2d %02d:%02d:%02d [%s] ", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), lg.logLevelToString(level))
}

func (log *Logger) logLevelToString(level LogLevel) string {
	switch level {
	case Critial:
		return "Critial"
	case Fatal:
		return "Fatal"
	case Error:
		return "Error"
	case Warn:
		return "Warn"
	case Info:
		return "Info"
	case Debug:
		return "Debug"
	}
	return "Unknown"
}
