package core

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/MCSManager/pty/utils"
	"github.com/mattn/go-colorable"
)

var PtySize string
var ColorAble bool
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
	// Remove the stdin cache, so that the system signal is passed directly to the PTY
	// This method operates on file descriptors and is not applicable to parent-child processes
	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	pty.resizeWindow(&PtySize)
	io.Copy(pty.StdIn(), utils.Encoder(Coder, os.Stdin))
}

func (pty *Pty) handleStdOut() {
	var stdout io.Writer
	if ColorAble {
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
