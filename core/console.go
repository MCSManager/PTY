//go:build !windows
// +build !windows

package console

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"

	"github.com/MCSManager/pty/core/interfaces"
)

var _ interfaces.Console = (*console)(nil)

type console struct {
	file  *os.File
	cmd   *exec.Cmd
	coder string

	initialCols uint
	initialRows uint

	env []string
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

func (c *console) stdIn() *os.File {
	return c.file
}

func (c *console) stdOut() *os.File {
	return c.file
}

func (c *console) SetSize(cols uint, rows uint) {
	c.initialRows = rows
	c.initialCols = cols
	if c.file == nil {
		return
	}
	pty.Setsize(c.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
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
