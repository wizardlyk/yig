package log

import (
	"fmt"
	"io"
	"log"
	"os"
	//	"github.com/journeymidnight/yig/helper"
)


type TracerLogger struct {
	Logger   *log.Logger
	LogLevel int
}

func NewTracerLog(out io.Writer, prefix string, flag int, level int) *TracerLogger {
	var logger TracerLogger
	logger.LogLevel = level
	logger.Logger = log.New(out, prefix, flag)
	return &logger
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *TracerLogger) Printf(level int, format string, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *TracerLogger) Print(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprint(v...))
	}
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *TracerLogger) Println(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintln(v...))
	}
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *TracerLogger) Fatal(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprint(v...))
	}
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *TracerLogger) Fatalf(level int, format string, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintf(format, v...))
	}
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *TracerLogger)Fatalln(level int, v ...interface{}) {
	if l.LogLevel >= level {
		l.Logger.Output(2, fmt.Sprintln(v...))
	}
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *TracerLogger) Panic(level int, v ...interface{}) {
	s := fmt.Sprint(v...)
	if l.LogLevel >= level {
		l.Logger.Output(2, s)
	}
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *TracerLogger) Panicf(level int, format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	if l.LogLevel >= level {
		l.Logger.Output(2, s)
	}
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *TracerLogger) Panicln(level int, v ...interface{}) {
	s := fmt.Sprintln(v...)
	if l.LogLevel >= level {
		l.Logger.Output(2, s)
	}
	panic(s)
}
