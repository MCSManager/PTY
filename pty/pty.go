//go:build !windows
// +build !windows

package pty

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	opty "github.com/creack/pty"
	"github.com/mattn/go-colorable"
)

type Pty struct {
	tty *os.File
	cmd *exec.Cmd
}

type Getdata struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func Start(dir, command string) (*Pty, error) {
	cmd := exec.Command(command)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "TERM=xterm")
	tty, err := opty.Start(cmd)
	return &Pty{tty: tty, cmd: cmd}, err
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

func (pty *Pty) HandleStd() {
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
	// go func() { _, _ = io.Copy(pty.tty, os.Stdin) }()
	// io.Copy(os.Stdout, pty.tty)
	go pty.handleStdOut()
	pty.handleStdIn()
}

func (pty *Pty) handleStdIn() {
	inputReader := bufio.NewReader(os.Stdin)
	var getdata Getdata
	for {
		input, _ := inputReader.ReadString('\n')
		json.Unmarshal([]byte(input), &getdata)
		switch getdata.Type {
		case 1:
			pty.tty.Write([]byte(getdata.Data + "\n"))
		case 2:
			cols, err := strconv.Atoi(strings.Split(getdata.Data, ",")[0])
			if err != nil {
				fmt.Println("SetSize err:", err)
				panic(err)
			}
			rows, err := strconv.Atoi(strings.Split(getdata.Data, ",")[1])
			if err != nil {
				fmt.Println("SetSize err:", err)
				panic(err)
			}
			pty.Setsize(uint32(cols), uint32(rows))
		case 3:
			pty.tty.Write([]byte{3})
		default:
			continue
		}
	}
}

func (pty *Pty) handleStdOut() {
	var err error
	var n int
	buf := make([]byte, 2*2048)
	reader := bufio.NewReader(pty.tty)
	stdout := colorable.NewColorableStdout()
	for {
		n, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Failed to read from pty master: %s", err)
			continue
		} else if err == io.EOF {
			pty.Close()
			os.Exit(0)
		}
		stdout.Write(buf[:n])
	}
}
