package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/MCSManager/pty/utils"
	"github.com/mattn/go-colorable"
)

var PtySize string
var Color bool
var Coder string

type DataProtocol struct {
	Type int    `json:"type"`
	Data string `json:"data"`
}

func (pty *Pty) HandleStdIO() {
	go pty.handleStdIn()
	pty.handleStdOut()
}

func (pty *Pty) handleStdIn() {
	if PtySize == "" {
		pty.Setsize(50, 50)
		pty.noSizeFlag()
	} else {
		pty.resizeWindow(&PtySize)
		pty.existSizeFlag()
	}
}

func (pty *Pty) noSizeFlag() {
	var err error
	var protocol DataProtocol
	var bufferText string
	var data []byte
	inputReader := bufio.NewReader(os.Stdin)
	for {
		bufferText, _ = inputReader.ReadString('\n')
		err = json.Unmarshal([]byte(bufferText), &protocol)
		if err != nil {
			fmt.Printf("[MCSMANAGER-PTY] Unmarshall json err: %v\noriginal data: %s\n", err, bufferText)
			continue
		}
		switch protocol.Type {
		case 1:
			data, err = ioutil.ReadAll(utils.Encoder(Coder, bytes.NewReader([]byte(protocol.Data))))
			if err != nil {
				continue
			}
			pty.StdIn().Write(data)
		case 2:
			pty.resizeWindow(&protocol.Data)
		case 3:
			pty.StdIn().Write([]byte{03})
		default:
		}
	}
}

func (pty *Pty) existSizeFlag() {
	// Remove the stdin cache, so that the system signal is passed directly to the PTY
	// This method operates on file descriptors and is not applicable to parent-child processes
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	io.Copy(pty.StdIn(), utils.Encoder(Coder, os.Stdin))
}

func (pty *Pty) handleStdOut() {
	var stdout io.Writer
	if Color {
		stdout = colorable.NewColorable(os.Stdout)
	} else {
		stdout = colorable.NewNonColorable(os.Stdout)
	}
	io.Copy(stdout, utils.Decoder(Coder, pty.StdOut()))
}

// Set the PTY window size based on the text
func (pty *Pty) resizeWindow(sizeText *string) {
	arr := strings.Split(*sizeText, ",")
	if len(arr) != 2 {
		fmt.Printf("[MCSMANAGER-PTY] Set PTY size data failed,original data:%#v\n", *sizeText)
		return
	}
	cols, err1 := strconv.Atoi(arr[0])
	rows, err2 := strconv.Atoi(arr[1])
	if err1 != nil || err2 != nil {
		fmt.Printf("[MCSMANAGER-PTY] Failed to set window size,original data:%#v\n", *sizeText)
		return
	}
	pty.Setsize(uint32(cols), uint32(rows))
}
