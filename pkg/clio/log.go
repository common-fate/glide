// Package clio contains helpers for printing CLI output messages.
package clio

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// Log prints to stdout.
func Log(format string, a ...interface{}) {
	color.White(format, a...)
}

// Info prints to stderr with an [i] indicator.
func Info(format string, a ...interface{}) {
	format = "[i] " + format
	fmt.Fprintln(color.Error, color.WhiteString(format, a...))
}

// Success prints to stderr with a [✔] indicator.
func Success(format string, a ...interface{}) {
	format = "[✔] " + format
	fmt.Fprintln(color.Error, color.GreenString(format, a...))
}

// Error prints to stderr with a [✘] indicator.
func Error(format string, a ...interface{}) {
	format = "[✘] " + format
	fmt.Fprintln(color.Error, color.RedString(format, a...))
}

// Warn prints to stderr with a [!] indicator.
func Warn(format string, a ...interface{}) {
	format = "[!] " + format
	fmt.Fprintln(color.Error, color.YellowString(format, a...))
}

// Warn prints to stderr with a [DEBUG] indicator
// if the GRANTED_LOG environment variable is set to 'debug'.
func Debug(format string, a ...interface{}) {
	if strings.ToLower(os.Getenv("GRANTED_LOG")) == "debug" {
		format = "[DEBUG] " + format
		fmt.Fprintln(color.Error, color.HiBlackString(format, a...))
	}
}
