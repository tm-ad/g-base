package log

import (
	"fmt"
	"github.com/tm-ad/g-base/util/pio"
	"io"
	"os"
	"sync"
	"time"
)

// Logger is our golog 简化版.
type Logger struct {
	Prefix     []byte
	Level      Level
	TimeFormat string
	// if new line should be added on all log functions, even the `F`s.
	// It defaults to true.
	//
	// See `golog#NewLine(newLineChar string)` as well.
	//
	// Note that this will not override the time and level prefix,
	// if you want to customize the log message please read the examples
	// or navigate to: https://github.com/kataras/golog/issues/3#issuecomment-355895870.
	NewLine bool
	mu      sync.Mutex
	Printer *pio.Printer
	// handlers []Handler
	once sync.Once
	logs sync.Pool
	// children *loggerMap
}

// New returns a new golog with a default output to `os.Stdout`
// and level to `InfoLevel`.
func New() *Logger {
	return &Logger{
		Level:      InfoLevel,
		TimeFormat: "2006/01/02 15:04",
		NewLine:    true,
		Printer:    pio.NewPrinter("", os.Stdout).EnableDirectOutput().Hijack(logHijacker),
		// children:   newLoggerMap(),
	}
}

// acquireLog returns a new log fom the pool.
func (l *Logger) acquireLog(level Level, msg string, withPrintln bool) *Log {
	log, ok := l.logs.Get().(*Log)
	if !ok {
		log = &Log{
			Logger: l,
		}
	}
	log.NewLine = withPrintln
	log.Time = time.Now()
	log.Level = level
	log.Message = msg
	return log
}

// releaseLog Log releases a log instance back to the pool.
func (l *Logger) releaseLog(log *Log) {
	l.logs.Put(log)
}

// we could use marshal inside Log but we don't have access to printer,
// we could also use the .Handle with NopOutput too but
// this way is faster:
var logHijacker = func(ctx *pio.Ctx) {
	l, ok := ctx.Value.(*Log)
	if !ok {
		ctx.Next()
		return
	}

	line := GetTextForLevel(l.Level, ctx.Printer.IsTerminal)
	if line != "" {
		line += " "
	}

	if t := l.FormatTime(); t != "" {
		line += t + " "
	}
	line += l.Message

	var b []byte
	if pref := l.Logger.Prefix; len(pref) > 0 {
		b = append(pref, []byte(line)...)
	} else {
		b = []byte(line)
	}

	ctx.Store(b, nil)
	ctx.Next()
}

// NopOutput disables the output.
var NopOutput = pio.NopOutput()

// SetOutput overrides the Logger's Printer's Output with another `io.Writer`.
//
// Returns itself.
func (l *Logger) SetOutput(w io.Writer) *Logger {
	l.Printer.SetOutput(w)
	return l
}

// AddOutput adds one or more `io.Writer` to the Logger's Printer.
//
// If one of the "writers" is not a terminal-based (i.e File)
// then colors will be disabled for all outputs.
//
// Returns itself.
func (l *Logger) AddOutput(writers ...io.Writer) *Logger {
	l.Printer.AddOutput(writers...)
	return l
}

// SetPrefix sets a prefix for this "l" Logger.
//
// The prefix is the first space-separated
// word that is being presented to the output.
// It's written even before the log level text.
//
// Returns itself.
func (l *Logger) SetPrefix(s string) *Logger {
	l.mu.Lock()
	l.Prefix = []byte(s)
	l.mu.Unlock()
	return l
}

// SetTimeFormat sets time format for logs,
// if "s" is empty then time representation will be off.
//
// Returns itself.
func (l *Logger) SetTimeFormat(s string) *Logger {
	l.mu.Lock()
	l.TimeFormat = s
	l.mu.Unlock()

	return l
}

// DisableNewLine disables the new line suffix on every log function, even the `F`'s,
// the caller should add "\n" to the log message manually after this call.
//
// Returns itself.
func (l *Logger) DisableNewLine() *Logger {
	l.mu.Lock()
	l.NewLine = false
	l.mu.Unlock()

	return l
}

// SetLevel accepts a string representation of
// a `Level` and returns a `Level` value based on that "levelName".
//
// Available level names are:
// "disable"
// "fatal"
// "error"
// "warn"
// "info"
// "debug"
//
// Alternatively you can use the exported `Level` field, i.e `Level = golog.ErrorLevel`
//
// Returns itself.
func (l *Logger) SetLevel(levelName string) *Logger {
	l.mu.Lock()
	l.Level = ParseLevel(levelName)
	l.mu.Unlock()

	return l
}

func (l *Logger) print(level Level, msg string, newLine bool) {
	if l.Level >= level {
		// newLine passed here in order for handler to know
		// if this message derives from Println and Leveled functions
		// or by simply, Print.
		log := l.acquireLog(level, msg, newLine)
		// if not handled by one of the handler
		// then print it as usual.
		// if !l.handled(log) {
		if newLine {
			l.Printer.Println(log)
		} else {
			l.Printer.Print(log)
		}
		// }

		l.releaseLog(log)
	}
	// if level was fatal we don't care about the logger's level, we'll exit.
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Print prints a log message without levels and colors.
func (l *Logger) Print(v ...interface{}) {
	l.print(DisableLevel, fmt.Sprint(v...), l.NewLine)
}

// Printf formats according to a format specifier and writes to `Printer#Output` without levels and colors.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.print(DisableLevel, fmt.Sprintf(format, args...), l.NewLine)
}

// Println prints a log message without levels and colors.
// It adds a new line at the end, it overrides the `NewLine` option.
func (l *Logger) Println(v ...interface{}) {
	l.print(DisableLevel, fmt.Sprint(v...), true)
}

// Log prints a leveled log message to the output.
// This method can be used to use custom log levels if needed.
// It adds a new line in the end.
func (l *Logger) Log(level Level, v ...interface{}) {
	l.print(level, fmt.Sprint(v...), l.NewLine)
}

// Logf prints a leveled log message to the output.
// This method can be used to use custom log levels if needed.
// It adds a new line in the end.
func (l *Logger) Logf(level Level, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Log(level, msg)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func (l *Logger) Fatal(v ...interface{}) {
	l.Log(FatalLevel, v...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Fatal(msg)
}

// Error will print only when logger's Level is error, warn, info or debug.
func (l *Logger) Error(v ...interface{}) {
	l.Log(ErrorLevel, v...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func (l *Logger) Errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Error(msg)
}

// Warn will print when logger's Level is warn, info or debug.
func (l *Logger) Warn(v ...interface{}) {
	l.Log(WarnLevel, v...)
}

// Warnf will print when logger's Level is warn, info or debug.
func (l *Logger) Warnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Warn(msg)
}

// Info will print when logger's Level is info or debug.
func (l *Logger) Info(v ...interface{}) {
	l.Log(InfoLevel, v...)
}

// Infof will print when logger's Level is info or debug.
func (l *Logger) Infof(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.Info(msg)
}

// Debug will print when logger's Level is debug.
func (l *Logger) Debug(v ...interface{}) {
	l.Log(DebugLevel, v...)
}

// Debugf will print when logger's Level is debug.
func (l *Logger) Debugf(format string, args ...interface{}) {
	// On debug mode don't even try to fmt.Sprintf if it's not required,
	// this can be used to allow `Debugf` to be called without even the `fmt.Sprintf`'s
	// performance cost if the logger doesn't allow debug logging.
	if l.Level >= DebugLevel {
		msg := fmt.Sprintf(format, args...)
		l.Debug(msg)
	}
}
