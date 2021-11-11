package logging

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logLevel Level
	out      *log.Logger
}

var logFile *os.File
var levelNames []string

func New() *Logger {
	logFile = os.Stdout
	levelNames = defaltLevelNames
	l := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile|log.Lmicroseconds)
	return &Logger{
		logLevel: DEBUG,
		out:      l,
	}
}

func (l *Logger) GetLogLevel() Level {
	return l.logLevel
}

func (l *Logger) SetLogLevel(logLevel Level) {
	if logLevel > ERROR || logLevel < DEBUG {
		return
	}
	l.logLevel = logLevel
}

func (l *Logger) SetLogFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New(fmt.Sprintf("error open file: %v", err))
	}
	l.out.SetOutput(f)
	return nil
}

func (l *Logger) EnableLogColor() {
	levelNames = colorLevelNames
}

func (l *Logger) Close() {
	logFile.Close()
}

// GetLevelName 获取级别名称
func (l *Logger) GetLevelName(level Level) string {
	return levelNames[int(level)]
}

// Output writes the output for a logging event
func (l *Logger) Output(calldepth int, s string) error {
	calldepth += 1
	return l.out.Output(calldepth, s)
}

func (l *Logger) Errorln(args ...interface{}) {
	if ERROR < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[ERROR]+fmt.Sprintln(args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	if ERROR < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[ERROR]+fmt.Sprintf(format, args...))
}

func (l *Logger) Warnln(args ...interface{}) {
	if WARN < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[WARN]+fmt.Sprintln(args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if WARN < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[WARN]+fmt.Sprintf(format, args...))
}

func (l *Logger) Infoln(args ...interface{}) {
	if INFO < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[INFO]+fmt.Sprintln(args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if INFO < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[INFO]+fmt.Sprintf(format, args...))
}

func (l *Logger) Debugln(args ...interface{}) {
	if DEBUG < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[DEBUG]+fmt.Sprintln(args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if DEBUG < l.logLevel {
		return
	}
	l.out.Output(2, levelNames[DEBUG]+fmt.Sprintf(format, args...))
}
