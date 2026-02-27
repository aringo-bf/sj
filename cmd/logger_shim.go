package cmd

import (
	"fmt"
	"strings"
)

// log provides a minimal compatibility layer for legacy logrus-style calls.
var log diagLogger

type diagLogger struct{}

func (diagLogger) Info(args ...interface{}) {
	printInfo("%s", fmt.Sprint(args...))
}

func (diagLogger) Infof(format string, args ...interface{}) {
	printInfo(format, args...)
}

func (diagLogger) Warn(args ...interface{}) {
	printWarn("%s", fmt.Sprint(args...))
}

func (diagLogger) Warnf(format string, args ...interface{}) {
	printWarn(format, args...)
}

func (diagLogger) Warnln(args ...interface{}) {
	msg := strings.TrimSuffix(fmt.Sprintln(args...), "\n")
	printWarn("%s", msg)
}

func (diagLogger) Error(args ...interface{}) {
	printErr("%s", fmt.Sprint(args...))
}

func (diagLogger) Errorf(format string, args ...interface{}) {
	printErr(format, args...)
}

func (diagLogger) Printf(format string, args ...interface{}) {
	printWarn(format, args...)
}

func (diagLogger) Fatal(args ...interface{}) {
	die("%s", fmt.Sprint(args...))
}

func (diagLogger) Fatalf(format string, args ...interface{}) {
	die(format, args...)
}
