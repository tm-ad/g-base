package log

import (
	"github.com/tm-ad/g-base/exceptions"
	"github.com/tm-ad/g-base/util"
	"time"
)

var _loggerInDevelopment *Logger

// A Log represents a log line.
type Log struct {
	// Logger is the original printer of this Log.
	Logger *Logger
	// Time is the current time
	Time time.Time
	// Level is the log level.
	Level Level
	// Message is the string reprensetation of the log's main body.
	Message string
	// NewLine returns false if this Log
	// derives from a `Print` function,
	// otherwise true if derives from a `Println`, `Error`, `Errorf`, `Warn`, etc...
	//
	// This NewLine does not mean that `Message` ends with "\n" (or `pio#NewLine`).
	// NewLine has to do with the methods called,
	// not the original content of the `Message`.
	NewLine bool
}

// FormatTime returns the formatted `Time`.
func (l Log) FormatTime() string {
	if l.Logger.TimeFormat == "" {
		return ""
	}
	return l.Time.Format(l.Logger.TimeFormat)
}

// TipInDevelopment 提供在开发模式下提示开发人员处理的信息
func TipInDevelopment(msg string) {
	if util.Development() {
		if _loggerInDevelopment == nil {
			_loggerInDevelopment = New()
		}
		_loggerInDevelopment.Warn(msg)
		_loggerInDevelopment.Warn(exceptions.CallStack())
	}
}
