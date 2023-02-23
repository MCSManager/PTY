package console

import (
	"bytes"
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
	mutex "github.com/juju/mutex/v2"
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
	r, err := mutex.Acquire(mutex.Spec{Name: "pty-winpty-lock", Timeout: time.Second * 5, Delay: time.Millisecond * 3, Clock: &fakeClock{}})
	if err != nil {
		return err
	}
	defer r.Release()
	dllDir, err := c.findDll(true)
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

	var pty *winpty.WinPTY
	if pty, err = winpty.OpenWithOptions(option); err != nil {
		if option.DllDir, err = c.findDll(false); err != nil {
			return err
		}
		if pty, err = winpty.OpenWithOptions(option); err != nil {
			return err
		}
	}
	c.stdIn = pty.Stdin
	c.stdOut = pty.Stdout
	c.stdErr = pty.Stderr
	c.file = pty
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

type fakeClock struct {
	delay time.Duration
}

func (f *fakeClock) After(time.Duration) <-chan time.Time {
	return time.After(f.delay)
}

func (f *fakeClock) Now() time.Time {
	return time.Now()
}

func (c *console) findDll(SkipExistFile bool) (string, error) {
	dllDir := filepath.Join(os.TempDir(), "pty_winpty")

	if err := os.MkdirAll(dllDir, os.ModePerm); err != nil {
		return "", err
	}
	if err := utils.UnzipWithFile(bytes.NewReader(winpty_zip), dllDir, utils.UnzipCfg{
		CoderTypes:    utils.T_UTF8,
		SkipExistFile: SkipExistFile,
	}); err != nil {
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
