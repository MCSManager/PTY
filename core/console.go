//go:build !windows
// +build !windows

package console

import (
	"os"
	"os/exec"

	"github.com/creack/pty"

	"github.com/MCSManager/pty/core/interfaces"
)

var _ interfaces.Console = (*console)(nil)

type console struct {
	file  *os.File
	cmd   *exec.Cmd
	coder string

	initialCols int
	initialRows int

	env []string
}

func newNative(coder string) Console {
	return &console{
		initialCols: 50,
		initialRows: 50,
		coder:       coder,

		file: nil,

		env: append(os.Environ(), "TERM=xterm-256color"),
	}
}

func (c *console) Start(dir string, command []string) error {
	cmd, err := c.buildCmd(command)
	if err != nil {
		return err
	}
	c.cmd = cmd
	cmd.Dir = dir
	cmd.Env = c.env
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	c.file = f
	return nil
}

func (c *console) buildCmd(args []string) (*exec.Cmd, error) {
	if len(args) < 1 {
		return nil, ErrInvalidCmd
	}
	cmd := exec.Command(args[0], args[1:]...)
	return cmd, nil
}

func (c *console) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.Read(b)
}

func (c *console) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.Write(b)
}

func (c *console) stdIn() *os.File {
	return c.file
}

func (c *console) stdOut() *os.File {
	return c.file
}

func (c *console) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	return c.file.Close()
}

func (c *console) SetSize(cols int, rows int) error {
	if c.file == nil {
		c.initialRows = rows
		c.initialCols = cols
		return nil
	}

	return pty.Setsize(c.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (c *console) GetSize() (int, int, error) {
	if c.file == nil {
		return c.initialCols, c.initialRows, nil
	}

	rows, cols, err := pty.Getsize(c.file)
	return cols, rows, err
}

func (c *console) AddENV(environ []string) error {
	c.env = append(os.Environ(), environ...)
	return nil
}

func (c *console) Pid() int {
	if c.cmd == nil {
		return 0
	}

	return c.cmd.Process.Pid
}

func (c *console) Wait() (*os.ProcessState, error) {
	if c.cmd == nil {
		return nil, ErrProcessNotStarted
	}

	return c.cmd.Process.Wait()
}

func (c *console) Kill() error {
	if c.cmd == nil {
		return ErrProcessNotStarted
	}

	return c.cmd.Process.Kill()
}

func (c *console) Signal(sig os.Signal) error {
	if c.cmd == nil {
		return ErrProcessNotStarted
	}

	return c.cmd.Process.Signal(sig)
}
