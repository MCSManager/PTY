//go:build !windows
// +build !windows

package core

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	opty "github.com/creack/pty"
)

type Pty struct {
	tty    *os.File
	cmd    *exec.Cmd
	StdIn  *os.File
	StdOut *os.File
}

func Start(dir string, command []string) (*Pty, error) {
	// 去除命令参数外层引号
	for k, v := range command {
		if v[:1] == `"` && v[len(v)-1:] == `"` {
			command[k] = v[1 : len(v)-1]
		}
	}
	fmt.Printf("[MCSMANAGER-TTY] Full command: %s\n", command)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	cmd.Env = os.Environ()
	tty, err := opty.Start(cmd)
	return &Pty{tty: tty, cmd: cmd, StdIn: tty, StdOut: tty}, err
}

func (pty *Pty) Write(p []byte) (n int, err error) {
	return pty.tty.Write(p)
}

func (pty *Pty) Read(p []byte) (n int, err error) {
	return pty.tty.Read(p)
}

func (pty *Pty) Setsize(cols, rows uint32) error {
	return opty.Setsize(pty.tty, &opty.Winsize{
		Cols: uint16(cols),
		Rows: uint16(rows),
	})
}

func (pty *Pty) killChildProcess(c *exec.Cmd) error {
	pgid, err := syscall.Getpgid(c.Process.Pid)
	if err != nil {
		// Fall-back on error. Kill the main process only.
		c.Process.Kill()
	}
	// Kill the whole process group.
	syscall.Kill(-pgid, syscall.SIGTERM)
	return c.Wait()
}

func (pty *Pty) Close() error {
	if err := pty.tty.Close(); err != nil {
		return err
	}
	return pty.killChildProcess(pty.cmd)
}
