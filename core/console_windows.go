package console

import (
	"archive/zip"
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/MCSManager/pty/core/go-winpty"
	"github.com/MCSManager/pty/core/interfaces"
	"github.com/MCSManager/pty/utils"
)

//go:embed winpty
var winpty_zip []byte

var _ interfaces.Console = (*console)(nil)

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

func (c *console) Start(dir string, command []string) error {
	dllDir, err := c.UnloadEmbeddedDeps()
	if err != nil {
		return err
	}

	option := winpty.Options{
		DllDir:      dllDir,
		Command:     c.buildCmd(command),
		Dir:         dir,
		Env:         c.env,
		InitialCols: uint32(c.initialCols),
		InitialRows: uint32(c.initialRows),
	}

	if c.colorAble {
		option.AgentFlags = winpty.WINPTY_FLAG_COLOR_ESCAPES
	} else {
		option.AgentFlags = winpty.WINPTY_FLAG_PLAIN_OUTPUT
	}
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

func (c *console) buildCmd(args []string) string {
	var cmds = fmt.Sprintf("cmd /C chcp %s > nul & ", utils.CodePage(c.coder))
	for _, v := range args {
		cmds += fmt.Sprintf("%s ", v)
	}
	return cmds
}

func (c *console) UnloadEmbeddedDeps() (string, error) {
	dllDir := filepath.Join(os.TempDir(), "pty_winpty")
	if err := os.MkdirAll(dllDir, 0755); err != nil {
		return "", err
	}

	dirInfo, err := os.ReadDir(dllDir)
	if err != nil {
		return "", err
	}
	if len(dirInfo) != 2 {
		unzip(bytes.NewReader(winpty_zip), dllDir)
	}
	return dllDir, nil
}

func unzip(f *bytes.Reader, targetPath string) error {
	zipReader, err := zip.NewReader(f, f.Size())
	if err != nil {
		return err
	}
	for _, f := range zipReader.File {
		fpath := filepath.Join(targetPath, f.Name)
		info, statErr := os.Stat(fpath)
		if statErr == nil && f.FileInfo().Size() == info.Size() {
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		inFile, err := f.Open()
		if err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, inFile)
		if err != nil {
			return err
		}
		inFile.Close()
		outFile.Close()
	}
	return err
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

func (c *console) SetSize(cols uint, rows uint) error {
	c.initialRows = rows
	c.initialCols = cols
	if c.file == nil {
		return nil
	}
	err := c.file.SetSize(uint32(c.initialCols), uint32(c.initialRows))
	if err.Error() != "The operation completed successfully." {
		return err
	}
	return nil
}

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

func (c *console) Kill() error {
	_, err := c.findProcess()
	if err != nil {
		return err
	}
	return exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(c.Pid())).Run()
}
