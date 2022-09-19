package console

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/MCSManager/pty/console/iface"
)

var (
	ErrProcessNotStarted = errors.New("process has not been started")
	ErrInvalidCmd        = errors.New("invalid command")
)

type Console iface.Console

// Create a new pty
func New(coder string, colorAble bool) Console {
	return newNative(coder, colorAble, 50, 50)
}

// Create a new pty and initialize the size
func NewWithSize(coder string, colorAble bool, Cols, Rows uint) Console {
	return newNative(coder, colorAble, Cols, Rows)
}

func newNative(coder string, colorAble bool, Cols, Rows uint) Console {
	if Cols == 0 {
		Cols = 50
	}
	if Rows == 0 {
		Rows = 50
	}
	console := console{
		initialCols: Cols,
		initialRows: Rows,
		coder:       coder,
		colorAble:   colorAble,

		file: nil,
	}
	if runtime.GOOS == "windows" {
		console.env = os.Environ()
	} else {
		console.env = append(os.Environ(), "TERM=xterm-256color")
	}
	return &console
}

// Read data from pty console
func (c *console) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.StdOut().Read(b)
}

// Write data to the pty console
func (c *console) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.StdIn().Write(b)
}

func (c *console) StdIn() io.Writer {
	return c.stdIn
}

func (c *console) StdOut() io.Reader {
	return c.stdOut
}

func (c *console) StdErr() io.Reader {
	return c.stdErr
}

// Add environment variables before start
func (c *console) AddENV(environ []string) error {
	c.env = append(c.env, environ...)
	return nil
}

// close the pty and kill the subroutine
func (c *console) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	return c.file.Close()
}

// wait for the pty subroutine to exit
func (c *console) Wait() (*os.ProcessState, error) {
	proc, err := c.findProcess()
	if err != nil {
		return nil, err
	}
	return proc.Wait()
}

// Send system signals to pty subroutines
func (c *console) Signal(sig os.Signal) error {
	proc, err := c.findProcess()
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}

// ResizeWithString("50,50")
func (c *console) ResizeWithString(sizeText string) error {
	arr := strings.Split(sizeText, ",")
	if len(arr) != 2 {
		return fmt.Errorf("the parameter is incorrect")
	}
	cols, err1 := strconv.Atoi(arr[0])
	rows, err2 := strconv.Atoi(arr[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("failed to set window size")
	}
	if cols < 0 || rows < 0 {
		return fmt.Errorf("failed to set window size")
	}
	return c.SetSize(uint(cols), uint(rows))
}

// Get pty window size
func (c *console) GetSize() (uint, uint) {
	return c.initialCols, c.initialRows
}
