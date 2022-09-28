package console

import (
	"archive/zip"
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/MCSManager/pty/console/go-winpty"
	"github.com/MCSManager/pty/console/iface"
	"github.com/juju/fslock"
)

//go:embed winpty
var winpty_zip []byte

var _ iface.Console = (*console)(nil)

type console struct {
	file      *winpty.WinPTY
	coder     string
	colorAble bool

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

	if c.colorAble {
		option.AgentFlags = winpty.WINPTY_FLAG_COLOR_ESCAPES
	} else {
		// Do not output escape characters
		option.AgentFlags = winpty.WINPTY_FLAG_PLAIN_OUTPUT
	}
	// creat stderr
	option.AgentFlags = option.AgentFlags | winpty.WINPTY_FLAG_CONERR
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
	var cmds = fmt.Sprintf("cmd /C chcp %s > nul & ", codePage(c.coder))
	if file, err := exec.LookPath(args[0]); err != nil {
		return "", err
	} else if args[0], err = filepath.Abs(file); err != nil {
		return "", err
	}
	for _, v := range args {
		cmds += v + ` `
	}
	return cmds[:len(cmds)-1], nil
}

var chcp = map[string]string{
	"UTF-8":     "65001",
	"UTF-16":    "1200",
	"GBK":       "936",
	"GB2312":    "936",
	"GB18030":   "54936",
	"BIG5":      "950",
	"KS_C_5601": "949",
	"SHIFTJIS":  "932",
}

func codePage(types string) string {
	if cp, ok := chcp[strings.ToUpper(types)]; ok {
		return cp
	} else {
		return "65001"
	}
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
	if err := releases(bytes.NewReader(winpty_zip), dllDir); err != nil {
		return "", err
	}
	return dllDir, nil
}

// release winpty prepend
func releases(f *bytes.Reader, targetPath string) error {
	zipReader, err := zip.NewReader(f, f.Size())
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		fpath := filepath.Join(targetPath, f.Name)
		info, statErr := os.Stat(fpath)
		// Check if the file is complete
		if statErr == nil && f.FileInfo().Size() == info.Size() {
			continue
		}
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		buf := bufio.NewWriter(outFile)
		if _, err = io.Copy(buf, inFile); err != nil {
			return err
		}
		if err := buf.Flush(); err != nil {
			return err
		}
		if err := inFile.Close(); err != nil {
			return err
		}
		if err := outFile.Close(); err != nil {
			return err
		}
	}
	return nil
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
