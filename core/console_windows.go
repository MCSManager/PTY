package console

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/MCSManager/pty/core/go-winpty"
	"github.com/MCSManager/pty/core/interfaces"
	"github.com/MCSManager/pty/utils"
)

//go:embed winpty/*
var winpty_deps embed.FS

var _ interfaces.Console = (*console)(nil)

type console struct {
	initialCols int
	initialRows int
	coder       string

	file *winpty.WinPTY

	env []string
}

func newNative(coder string) Console {
	return &console{
		initialCols: 50,
		initialRows: 50,
		coder:       coder,

		file: nil,

		env: os.Environ(),
	}
}

func (c *console) Start(dir string, command []string) error {
	dllDir, err := c.UnloadEmbeddedDeps()
	if err != nil {
		return err
	}
	if err = os.Chdir(dir); err != nil {
		return err
	}
	opts := winpty.Options{
		DLLPrefix: dllDir,
		Command:   c.buildCmd(command),
		Dir:       dir,
		Env:       c.env,
	}

	cmd, err := winpty.OpenWithOptions(opts)
	if err != nil {
		return err
	}

	c.file = cmd
	return nil
}

func (c *console) buildCmd(args []string) string {
	var cmds = fmt.Sprintf("cmd /C chcp %s && ", utils.CodePage(c.coder))
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
		filenameEmbedded := fmt.Sprintf("winpty/%s", file)
		filenameDisk := path.Join(dllDir, file)

		_, statErr := os.Stat(filenameDisk)
		if statErr == nil {
			continue
		}

		data, err := winpty_deps.ReadFile(filenameEmbedded)
		if err != nil {
			return "", err
		}
		file, err := os.OpenFile(filenameDisk, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return "", err
		}
		if _, err := file.Write(data); err != nil {
			return "", err
		}
		file.Close()
	}

	return dllDir, nil
}

func (c *console) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	n, err := c.file.StdOut.Read(b)

	return n, err
}

func (c *console) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.StdIn.Write(b)
}

func (c *console) stdIn() *os.File {
	return c.file.StdIn
}

func (c *console) stdOut() *os.File {
	return c.file.StdOut
}

func (c *console) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	c.file.Close()
	return nil
}

func (c *console) SetSize(cols int, rows int) error {
	c.initialRows = rows
	c.initialCols = cols

	if c.file == nil {
		return nil
	}

	c.file.SetSize(uint32(c.initialCols), uint32(c.initialRows))
	return nil
}

func (c *console) GetSize() (int, int, error) {
	return c.initialCols, c.initialRows, nil
}

func (c *console) AddENV(environ []string) error {
	c.env = append(c.env, environ...)
	return nil
}

func (c *console) Pid() int {
	return c.file.GetPid()
}

func (c *console) Wait() (*os.ProcessState, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}

	proc, err := os.FindProcess(int(c.Pid()))
	if err != nil {
		return nil, err
	}
	return proc.Wait()
}

func (c *console) Kill() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	proc, err := os.FindProcess(int(c.Pid()))
	if err != nil {
		return err
	}

	return proc.Kill()
}

func (c *console) Signal(sig os.Signal) error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	proc, err := os.FindProcess(int(c.Pid()))
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}
