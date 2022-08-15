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
