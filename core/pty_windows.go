//go:build windows

package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MCSManager/pty/core/winpty"
)

type Pty struct {
	tty *winpty.WinPTY
}

func Start(dir string, command []string) (*Pty, error) {
	path, err := getExecutableFilePath()
	if err != nil {
		return nil, err
	}
	var cmd string
	for _, v := range command {
		cmd += fmt.Sprintf("%s ", v)
	}
	if err = os.Chdir(dir); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Failed to change workspace directory: %v\n", err)
		panic(err)
	}
	tty, err := winpty.OpenWithOptions(winpty.Options{
		DLLPrefix: path,
		Command:   cmd,
		Dir:       dir,
		Env:       os.Environ(),
	})
	return &Pty{tty: tty}, err
}

func getExecutableFilePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	if executableFileExists(filepath.Dir(ex)+"/winpty-agent.exe") && executableFileExists(filepath.Dir(ex)+"/winpty.dll") {
		return filepath.Dir(ex), nil
	} else {
		return filepath.Dir(ex), errors.New("[MCSMANAGER-PTY] Executable File {winpty-agent.exe,winpty.dll} does not exist")
	}
}

func executableFileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func (pty *Pty) Write(p []byte) (n int, err error) {
	return pty.tty.StdIn.Write(p)
}

func (pty *Pty) Read(p []byte) (n int, err error) {
	return pty.tty.StdOut.Read(p)
}

func (pty *Pty) Setsize(cols, rows uint32) error {
	pty.tty.SetSize(cols, rows)
	return nil
}

func (pty *Pty) Close() error {
	pty.tty.Close()
	return nil
}

func (pty *Pty) StdOut() *os.File {
	return pty.tty.StdOut
}

func (pty *Pty) StdIn() *os.File {
	return pty.tty.StdIn
}

func (pty *Pty) Pid() int {
	return pty.tty.GetPid()
}
