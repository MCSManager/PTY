package console

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/MCSManager/pty/console/go-winpty"
	"github.com/MCSManager/pty/console/iface"
	"github.com/MCSManager/pty/utils"
	"github.com/juju/fslock"
)

//go:embed winpty
var winpty_zip []byte

var _ iface.Console = (*console)(nil)

type console struct {
	file  *winpty.WinPTY
	coder utils.CoderType

	stdIn  io.Writer
	stdOut io.Reader
	stdErr io.Reader

	initialCols uint
	initialRows uint

	env []string
}

// start pty subroutine
func (c *console) Start(dir string, command []string) error {
	dllDir, err := c.findDll()
	if err != nil {
		return err
	}
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	} else if err := os.Chdir(dir); err != nil {
		return err
	}
	cmd, err := c.buildCmd(command)
	if err != nil {
		return err
	}
	option := winpty.Options{
		DllDir:      dllDir,
		Command:     cmd,
		Dir:         dir,
		Env:         c.env,
		InitialCols: uint32(c.initialCols),
		InitialRows: uint32(c.initialRows),
	}

	// creat stderr
	option.AgentFlags = winpty.WINPTY_FLAG_CONERR | winpty.WINPTY_FLAG_COLOR_ESCAPES
	if cmd, err := winpty.OpenWithOptions(option); err != nil {
		return err
	} else {
		c.stdIn = cmd.Stdin
		c.stdOut = cmd.Stdout
		c.stdErr = cmd.Stderr
		c.file = cmd
	}
	return nil
}

// splice command
func (c *console) buildCmd(args []string) (string, error) {
	if len(args) == 0 {
		return "", ErrInvalidCmd
	}
	var cmds = fmt.Sprintf("cmd /C chcp %s > nul & ", utils.CodePage(c.coder))
	for _, v := range args {
		cmds += v + ` `
	}
	return cmds[:len(cmds)-1], nil
}

func (c *console) findDll() (string, error) {
	// File locks prevent concurrent file corruption
	flock := fslock.New(filepath.Join(os.TempDir(), "pty_winpty_lock"))
	if err := flock.LockWithTimeout(time.Second * 5); err != nil {
		return "", err
	}
	defer flock.Unlock()

	dllDir := filepath.Join(os.TempDir(), "pty_winpty")

	if err := os.MkdirAll(dllDir, os.ModePerm); err != nil {
		return "", err
	}
	if err := utils.UnzipWithFile(winpty_zip, dllDir, utils.T_UTF8); err != nil {
		return "", err
	}
	return dllDir, nil
}

// set pty window size
func (c *console) SetSize(cols uint, rows uint) error {
	c.initialRows = rows
	c.initialCols = cols
	if c.file == nil {
		return nil
	}
	err := c.file.SetSize(uint32(c.initialCols), uint32(c.initialRows))
	// Error special handling
	if err.Error() != "The operation completed successfully." {
		return err
	}
	return nil
}

// Get the process id of the pty subprogram
func (c *console) Pid() int {
	if c.file == nil {
		return 0
	}

	return c.file.Pid()
}

func (c *console) findProcess() (*os.Process, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}
	return os.FindProcess(c.Pid())
}

// Force kill pty subroutine
func (c *console) Kill() error {
	_, err := c.findProcess()
	if err != nil {
		return err
	}
	// try to kill all child processes
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(c.Pid())).Run()
}
