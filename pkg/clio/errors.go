package clio

import "fmt"

type Printer interface {
	Print()
}

type PrintCLIErrorer interface {
	PrintCLIError()
}
type CLIError struct {
	// Err is a string to avoid puntuation linting
	Err string
	// setting ExcludeDefaultError to true will cause the PrintCLIError method to only print the contents of Messages
	// by default PrintCLIError first prints the contents of Err, followed by the contents of Messages
	ExcludeDefaultError bool
	// Messages are items which implement the clio.Printer interface.
	// when PrintCLIError is called, each of the messages Print() method is called in order of appearence in the slice
	Messages []Printer
}

func NewCLIError(err string, msgs ...Printer) *CLIError {
	return &CLIError{
		Err:      err,
		Messages: msgs,
	}
}

type MsgType string

const (
	LogMsgType     MsgType = "LOG"
	InfoMsgType    MsgType = "INFO"
	SuccessMsgType MsgType = "SUCCESS"
	ErrorMsgType   MsgType = "ERROR"
	WarnMsgType    MsgType = "WARN"
	DebugMsgType   MsgType = "DEBUG"
)

type Msg struct {
	Msg  string
	Type MsgType
}

func (m Msg) Print() {
	switch m.Type {
	case InfoMsgType:
		Info(m.Msg)
	case SuccessMsgType:
		Success(m.Msg)
	case ErrorMsgType:
		Error(m.Msg)
	case WarnMsgType:
		Warn(m.Msg)
	case DebugMsgType:
		Debug(m.Msg)

	default:
		Log(m.Msg)
	}
}
func LogMsg(format string, a ...interface{}) Msg {
	return Msg{Msg: fmt.Sprintf(format, a...), Type: LogMsgType}
}
func InfoMsg(format string, a ...interface{}) Msg {
	return Msg{Msg: fmt.Sprintf(format, a...), Type: InfoMsgType}
}
func SuccessMsg(format string, a ...interface{}) Msg {
	return Msg{Msg: fmt.Sprintf(format, a...), Type: SuccessMsgType}
}
func ErrorMsg(format string, a ...interface{}) Msg {
	return Msg{Msg: fmt.Sprintf(format, a...), Type: ErrorMsgType}
}
func WarnMsg(format string, a ...interface{}) Msg {
	return Msg{Msg: fmt.Sprintf(format, a...), Type: WarnMsgType}
}
func DebugMsg(format string, a ...interface{}) Msg {
	return Msg{Msg: fmt.Sprintf(format, a...), Type: DebugMsgType}
}

// Error implements the error interface. It uses the default message of the
// wrapped error.
func (e *CLIError) Error() string {
	return e.Err
}

// PrintCLIError prints the error message and then any messages in order from the slice
// The indended use is to surface errors with useful messages then os.Exit without having to place os.Exit withing methods other than the cli main function
//
// err := CLIError{Err: errors.New("new error"), Messages: []Printer{LogMsg("hello world")}}
//
// err.PrintCLIError()
// // produces
// [e] new error
// hello world
func (e *CLIError) PrintCLIError() {
	if !e.ExcludeDefaultError {
		Error("%s", e.Err)
	}

	for i := range e.Messages {
		e.Messages[i].Print()
	}
}
