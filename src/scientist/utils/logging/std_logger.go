package logging

import "fmt"

var std = New()

// EnableLogColor 开启颜色
// TODO: 这里不是线程安全。
func EnableColor() {
	levelNames = colorLevelNames
}

func GetStdLogger() *Logger {
	return std
}

func SetLogLevel(logLevel Level) {
	std.SetLogLevel(logLevel)
}

func SetLogFile(filename string) error {
	return std.SetLogFile(filename)
}

func Errorln(args ...interface{}) {
	if ERROR < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[ERROR]+fmt.Sprintln(args...))
}

func Errorf(format string, args ...interface{}) {
	if ERROR < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[ERROR]+fmt.Sprintf(format, args...))
}

func Warnln(args ...interface{}) {
	if WARN < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[WARN]+fmt.Sprintln(args...))
}

func Warnf(format string, args ...interface{}) {
	if WARN < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[WARN]+fmt.Sprintf(format, args...))
}

func Infoln(args ...interface{}) {
	if INFO < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[INFO]+fmt.Sprintln(args...))
}

func Infof(format string, args ...interface{}) {
	if INFO < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[INFO]+fmt.Sprintf(format, args...))
}

func Debugln(args ...interface{}) {
	if DEBUG < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[DEBUG]+fmt.Sprintln(args...))
}

func Debugf(format string, args ...interface{}) {
	if DEBUG < std.logLevel {
		return
	}
	std.out.Output(2, levelNames[DEBUG]+fmt.Sprintf(format, args...))
}
