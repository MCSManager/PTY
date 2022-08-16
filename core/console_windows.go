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

//go:embed winpty/pty.zip
var winpty_zip []byte

var _ interfaces.Console = (*console)(nil)

type console struct {
	file  *winpty.WinPTY
	coder string

	initialCols int
	initialRows int

	env []string
}

func (c *console) Start(dir string, command []string) error {
	dllDir, err := c.UnloadEmbeddedDeps()
	if err != nil {
		return err
	}

	if cmd, err := winpty.OpenWithOptions(winpty.Options{
		DLLPrefix: dllDir,
		Command:   c.buildCmd(command),
		Dir:       dir,
		Env:       c.env,
	}); err != nil {
		return err
	} else {
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

	files := []string{"winpty.dll", "winpty-agent.exe"}
	for _, file := range files {
		_, statErr := os.Stat(filepath.Join(dllDir, file))
		if statErr == nil {
			continue
		} else {
			unzip(bytes.NewReader(winpty_zip), file, dllDir)
		}
	}
	return dllDir, nil
}

func unzip(f *bytes.Reader, fileName, targetPath string) error {
	zipReader, err := zip.NewReader(f, f.Size())
	if err != nil {
		return err
	}
	for _, f := range zipReader.File {
		if f.Name != fileName {
			continue
		}
		fpath := filepath.Join(targetPath, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
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
	}
	return err
}

func (c *console) stdIn() *os.File {
	return c.file.StdIn
}

func (c *console) stdOut() *os.File {
	return c.file.StdOut
}

func (c *console) SetSize(cols int, rows int) error {
	c.initialRows = rows
	c.initialCols = cols

	if c.file == nil {
		return nil
	}

	return c.file.SetSize(uint32(c.initialCols), uint32(c.initialRows))
}

func (c *console) Pid() int {
	if c.file == nil {
		return 0
	}
	return c.file.GetPid()
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
