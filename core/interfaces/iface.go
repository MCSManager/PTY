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

	SetSize(cols uint, rows uint) error

	ResizeWithString(sizeText string) error

	GetSize() (uint, uint)

	Start(dir string, command []string) error

	AddENV(environ []string) error

	Pid() int

	Wait() error

	Kill() error

	Signal(sig os.Signal) error

	HandleStdIO(ColorAble bool)
}
