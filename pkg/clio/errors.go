package clio

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

type LogMsg string

func (m LogMsg) Print() {
	Log(string(m))
}

type InfoMsg string

func (m InfoMsg) Print() {
	Info(string(m))
}

type WarnMsg string

func (m WarnMsg) Print() {
	Warn(string(m))
}

type DebugMsg string

func (m DebugMsg) Print() {
	Debug(string(m))
}

// Error implements the error interface. It uses the default message of the
// wrapped error.
func (e *CLIError) Error() string {
	return e.Err
}

// PrintCLIError prints the error message and then any messages in order from the slice
// The indended use is to surface errors with useful messages then os.Exit without having to place os.Exit withing methods other than the cli main function
//
// err := CLIError{Err: errors.New("new error"), Messages: []Printer{&LogMsg{Msg:"hello world"}}}
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
