package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

//
type LEVEL int

//
const (
	TRACE LEVEL = iota
	INFO
	WARN
	ERROR
	FATAL
)

//
var logger *log.Logger

var osExit = os.Exit
var formats = map[LEVEL]string{
	TRACE: "[TRACE] ",
	INFO:  "[ INFO] ",
	WARN:  "[ WARN] ",
	ERROR: "[ERROR] ",
	FATAL: "[FATAL] ",
}

// Sort: Trace, Info, Warn, Error, Fatal
var outputColors = []func(_ ...interface{}) string{
	color.New(color.FgMagenta).SprintFunc(),
	color.New(color.FgGreen).SprintFunc(),
	color.New(color.FgYellow).SprintFunc(),
	color.New(color.FgRed).SprintFunc(),
	color.New(color.FgHiRed).SprintFunc(),
}

func init() {
	logger = log.New(color.Output, "", log.Flags()&^(log.Ldate|log.Ltime))
}

// FormatMessage will return styled message
func FormatMessage(level LEVEL, format string, values ...interface{}) string {
	return formats[level] + fmt.Sprintf(format, values...)
}

// Write is shared output method
func Write(level LEVEL, format string, values ...interface{}) {
	message := FormatMessage(level, format, values...)

	logger.Print(outputColors[level](message))
}

// Trace for trace log
func Trace(format string, values ...interface{}) {
	Write(TRACE, format, values...)
}

// Info for info log
func Info(format string, values ...interface{}) {
	Write(INFO, format, values...)
}

// Warn for warn log
func Warn(format string, values ...interface{}) {
	Write(WARN, format, values...)
}

// Error for error log
func Error(format string, values ...interface{}) {
	Write(ERROR, format, values...)
}

// Fatal for fatal log
func Fatal(format string, values ...interface{}) {
	Write(FATAL, format, values...)

	osExit(1)
}
