//go:build windows

package pty

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/iamacarpet/go-winpty"
	"github.com/mattn/go-colorable"
)

type Pty struct {
	tty *winpty.WinPTY
}

type Getdata struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func getExecutableFilePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}

func Start(dir, command string) (*Pty, error) {
	path, err := getExecutableFilePath()
	if err != nil {
		return nil, err
	}
	tty, err := winpty.OpenWithOptions(winpty.Options{
		DLLPrefix: path,
		Command:   command,
		Dir:       dir,
		Env:       os.Environ(),
	})
	return &Pty{tty: tty}, err
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

func (pty *Pty) HandleStd() {
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
			pty.tty.StdIn.Write([]byte(getdata.Data))
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
			pty.tty.StdIn.Write([]byte{3})
		default:
			continue
		}
	}
}

func (pty *Pty) handleStdOut() {
	var err error
	var n int
	buf := make([]byte, 2*2048)
	reader := bufio.NewReader(pty.tty.StdOut)
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
