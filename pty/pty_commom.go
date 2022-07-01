package pty

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-colorable"
)

func (pty *Pty) HandleStdIO() {
	go pty.handleStdOut()
	pty.handleStdIn()
}

func (pty *Pty) handleStdIn() {
	inputReader := bufio.NewReader(os.Stdin)
	for {
		var err error
		var protocol DataProtocol
		bufferText, _ := inputReader.ReadString('\n')
		err = json.Unmarshal([]byte(bufferText), &protocol)
		if err != nil {
			fmt.Printf("[MCSMANAGER-TTY] Unmarshall json err:%v\n,original data:%#v\n", err, bufferText)
			continue
		}
		switch protocol.Type {
		case 1:
			pty.StdIn.Write([]byte(protocol.Data))
			continue
		case 2:
			resizeWindow(pty, protocol.Data)
			continue
		case 3:
			pty.StdIn.Write([]byte{3})
			continue
		default:
			continue
		}
	}
}

func (pty *Pty) handleStdOut() {
	var err error
	var n int
	buf := make([]byte, 2*2048)
	reader := bufio.NewReader(pty.StdOut)
	stdout := colorable.NewColorableStdout()
	for {
		n, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("[MCSMANAGER-TTY] Failed to read from pty master: %s", err)
			continue
		} else if err == io.EOF {
			pty.Close()
			os.Exit(-1)
		}
		stdout.Write(buf[:n])
	}
}

// Set the PTY window size based on the text
func resizeWindow(pty *Pty, sizeText string) {
	arr := strings.Split(sizeText, ",")
	cols, err1 := strconv.Atoi(arr[0])
	rows, err2 := strconv.Atoi(arr[0])
	if err1 != nil || err2 != nil {
		log.Printf("[MCSMANAGER-TTY] Failed to set window size: %s", err1)
		return
	}
	pty.Setsize(uint32(cols), uint32(rows))
}
