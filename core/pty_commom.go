package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-colorable"
)

type DataProtocol struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func (pty *Pty) HandleStdIO() {
	go pty.handleStdIn()
	pty.handleStdOut()
}

func (pty *Pty) handleStdIn() {
	var err error
	var protocol DataProtocol
	var bufferText string
	inputReader := bufio.NewReader(os.Stdin)
	for {
		bufferText, err = inputReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Printf("[MCSMANAGER-TTY] ReadString err: %v", err)
			continue
		}
		err = json.Unmarshal([]byte(bufferText), &protocol)
		if err != nil {
			fmt.Printf("[MCSMANAGER-TTY] Unmarshall json err: %v\n,original data: %#v\n", err, bufferText)
			continue
		}
		switch protocol.Type {
		case 1:
			pty.StdIn.Write([]byte(protocol.Data))
		case 2:
			resizeWindow(pty, protocol.Data)
		case 3:
			pty.StdIn.Write([]byte{3})
		default:
		}
	}
}

func (pty *Pty) handleStdOut() {
	var err error
	stdout := colorable.NewColorableStdout()
	_, err = io.Copy(stdout, pty.StdOut)
	if err != nil {
		fmt.Printf("[MCSMANAGER-TTY] Failed to read from pty master: %v\n", err)
		return
	}
}

// Set the PTY window size based on the text
func resizeWindow(pty *Pty, sizeText string) {
	arr := strings.Split(sizeText, ",")
	if len(arr) != 2 {
		fmt.Printf("[MCSMANAGER-TTY] Set tty size data failed,original data:%#v\n", sizeText)
		return
	}
	cols, err1 := strconv.Atoi(arr[0])
	rows, err2 := strconv.Atoi(arr[1])
	if err1 != nil || err2 != nil {
		fmt.Printf("[MCSMANAGER-TTY] Failed to set window size,original data:%#v\n", sizeText)
		return
	}
	pty.Setsize(uint32(cols), uint32(rows))
}
