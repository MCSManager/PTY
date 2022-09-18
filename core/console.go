//go:build !windows
// +build !windows

package console

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/creack/pty"

	"github.com/MCSManager/pty/core/interfaces"
)

var _ interfaces.Console = (*console)(nil)

type console struct {
	file      *os.File
	cmd       *exec.Cmd
	coder     string
	colorAble bool

	stdIn  io.Writer
	stdOut io.Reader
	stdErr io.Reader // nil

	initialCols uint
	initialRows uint

	env []string
}

func (c *console) Start(dir string, command []string) error {
	cmd, err := c.buildCmd(command)
	if err != nil {
		return err
	}
	if cwd, err := filepath.Abs(dir); err != nil {
		return err
	} else if err := os.Chdir(cwd); err != nil {
		return err
	}
	c.cmd = cmd
	cmd.Dir = dir
	cmd.Env = c.env
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: uint16(c.initialRows), Cols: uint16(c.initialCols)})
	if err != nil {
		return err
	}
	c.stdIn = DecoderWriter(c.coder, f)
	c.stdOut = DecoderReader(c.coder, f)
	c.stdErr = nil
	c.file = f
	return nil
}

func (c *console) buildCmd(args []string) (*exec.Cmd, error) {
	if len(args) < 1 {
		return nil, ErrInvalidCmd
	}
	if file, err := exec.LookPath(args[0]); err == nil {
		if path, err := filepath.Abs(file); err == nil {
			args[0] = path
		}
	}
	cmd := exec.Command(args[0], args[1:]...)
	return cmd, nil
}

func (c *console) StdIn() io.Writer {
	return c.stdIn
}

func (c *console) StdOut() io.Reader {
	return c.stdOut
}

func (c *console) StdErr() io.Reader {
	return nil
}

func (c *console) SetSize(cols uint, rows uint) error {
	c.initialRows = rows
	c.initialCols = cols
	if c.file == nil {
		return nil
	}
	return pty.Setsize(c.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (c *console) Pid() int {
	if c.cmd == nil {
		return 0
	}

	return c.cmd.Process.Pid
}

func (c *console) findProcess() (*os.Process, error) {
	if c.cmd == nil {
		return nil, ErrProcessNotStarted
	}
	return c.cmd.Process, nil
}

func (c *console) Kill() error {
	proc, err := c.findProcess()
	if err != nil {
		return err
	}
	pgid, err := syscall.Getpgid(proc.Pid)
	if err != nil {
		return proc.Kill()
	}
	return syscall.Kill(-pgid, syscall.SIGKILL)
}
