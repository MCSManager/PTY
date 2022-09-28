package iface

import (
	"io"
	"os"
)

// Console communication interface
type Console interface {
	io.Reader
	io.Writer
	// close the pty and kill the subroutine
	io.Closer

	// start pty subroutine
	Start(dir string, command []string) error

	// set pty window size
	SetSize(cols uint, rows uint) error

	// ResizeWithString("50,50")
	ResizeWithString(sizeText string) error

	// Get pty window size
	GetSize() (uint, uint)

	// Add environment variables before start
	AddENV(environ []string) error

	// Get the process id of the pty subprogram
	Pid() int

	// wait for the pty subroutine to exit
	Wait() (*os.ProcessState, error)

	// Force kill pty subroutine,try to kill all child processes
	Kill() error

	// Send system signals to pty subroutines
	Signal(sig os.Signal) error

	StdIn() io.Writer

	StdOut() io.Reader

	// nil in unix
	StdErr() io.Reader
}
