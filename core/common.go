package console

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/MCSManager/pty/core/interfaces"
	"github.com/MCSManager/pty/utils"
	"github.com/mattn/go-colorable"
)

var (
	ErrProcessNotStarted = errors.New("[MCSMANAGER-PTY] Process has not been started")
	ErrInvalidCmd        = errors.New("[MCSMANAGER-PTY] Invalid command")
)

type Console interfaces.Console

func New(coder string) Console {
	return newNative(coder)
}

func newNative(coder string) Console {
	console := console{
		initialCols: 50,
		initialRows: 50,
		coder:       coder,

		file: nil,
	}
	if runtime.GOOS == "windows" {
		console.env = os.Environ()
	} else {
		console.env = append(os.Environ(), "TERM=xterm-256color")
	}
	return &console
}

func (c *console) HandleStdIO(ColorAble bool) {
	go c.handleStdIn()
	go c.handleStdOut(ColorAble)
}

func (c *console) handleStdIn() {
	if runtime.GOOS == "windows" {
		io.Copy(c.stdIn(), os.Stdin)
	} else {
		io.Copy(c.stdIn(), utils.EncoderReader(c.coder, os.Stdin))
	}
}

func (c *console) handleStdOut(ColorAble bool) {
	var stdout io.Writer
	if ColorAble {
		stdout = colorable.NewColorable(os.Stdout)
	} else {
		stdout = colorable.NewNonColorable(os.Stdout)
	}
	if runtime.GOOS == "windows" {
		io.Copy(stdout, c.stdOut())
	} else {
		io.Copy(stdout, utils.DecoderReader(c.coder, c.stdOut()))
	}
}

func (c *console) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.stdOut().Read(b)
}

func (c *console) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.stdIn().Write(b)
}

func (c *console) AddENV(environ []string) error {
	c.env = append(c.env, environ...)
	return nil
}

func (c *console) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	return c.file.Close()
}

func (c *console) Wait() (*os.ProcessState, error) {
	proc, err := c.findProcess()
	if err != nil {
		return nil, err
	}
	return proc.Wait()
}

func (c *console) Kill() error {
	proc, err := c.findProcess()
	if err != nil {
		return err
	}

	return proc.Kill()
}

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
		return fmt.Errorf("[MCSMANAGER-PTY] The parameter is incorrect")
	}
	cols, err1 := strconv.Atoi(arr[0])
	rows, err2 := strconv.Atoi(arr[1])
	if err1 != nil || err2 != nil {
		return fmt.Errorf("[MCSMANAGER-PTY] Failed to set window size")
	}
	return c.SetSize(cols, rows)
}

func (c *console) GetSize() (int, int) {
	return c.initialCols, c.initialRows
}
