package diagnostics

import "fmt"

type Level string

const (
	InfoLevel  Level = "INFO"
	ErrorLevel Level = "ERROR"
	WarnLevel  Level = "WARN"
)

type Log struct {
	Level Level
	Msg   string
}

type Logs []Log

// HasSucceeded returns the success of the diagnostics.
// If there are any log entries with an ERROR level, it returns false.
// Otherwise, it returns true.
func (l *Logs) HasSucceeded() bool {
	for _, log := range *l {
		if log.Level == ErrorLevel {
			return false
		}
	}
	return true
}

func (l *Logs) Info(format string, a ...interface{}) {
	*l = append(*l, Log{
		Level: InfoLevel,
		Msg:   fmt.Sprintf(format, a...),
	})
}

func (l *Logs) Error(err error) {
	*l = append(*l, Log{
		Level: ErrorLevel,
		Msg:   err.Error(),
	})
}

func Info(format string, a ...interface{}) Logs {
	l := Logs{}
	l.Info(format, a...)
	return l
}

func Error(err error) Logs {
	l := Logs{}
	l.Error(err)
	return l
}
