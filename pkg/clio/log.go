// Package clio contains helpers for printing CLI output messages.
package clio

import (
	"os"
	"strings"

	"github.com/fatih/color"
)

func Log(format string, a ...interface{}) {
	color.White(format, a...)
}

func Info(format string, a ...interface{}) {
	format = "[i] " + format
	color.White(format, a...)
}

func Success(format string, a ...interface{}) {
	format = "[✔] " + format
	color.Green(format, a...)
}

func Error(format string, a ...interface{}) {
	format = "[✘] " + format
	color.Red(format, a...)
}

func Warn(format string, a ...interface{}) {
	format = "[!] " + format
	color.Yellow(format, a...)
}

func Debug(format string, a ...interface{}) {
	if strings.ToLower(os.Getenv("GRANTED_LOG")) == "debug" {
		format = "[DEBUG] " + format
		color.HiBlack(format, a...)
	}
}
