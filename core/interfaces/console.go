package interfaces

import (
	"io"
	"os"
)

// Console communication interface
type Console interface {
	io.Reader
	io.Writer
	io.Closer

	SetSize(cols int, rows int) error

	ResizeWithString(sizeText string) error

	GetSize() (int, int)

	Start(dir string, command []string) error

	AddENV(environ []string) error

	Pid() int

	Wait() (*os.ProcessState, error)

	Kill() error

	Signal(sig os.Signal) error

	HandleStdIO(ColorAble bool)
}
