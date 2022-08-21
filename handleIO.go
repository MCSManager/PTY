package main

import (
	"io"
	"os"
	"runtime"

	pty "github.com/MCSManager/pty/core"
	"github.com/MCSManager/pty/utils"
	"github.com/mattn/go-colorable"
)

func HandleStdIO(c pty.Console) {
	go handleStdIn(c)
	go handleStdOut(c)
	go handleStdErr(c)
}

func handleStdIn(c pty.Console) {
	if runtime.GOOS == "windows" {
		io.Copy(c.StdIn(), os.Stdin)
	} else {
		io.Copy(c.StdIn(), utils.EncoderReader(coder, os.Stdin))
	}
}

func handleStdOut(c pty.Console) {
	var stdout io.Writer
	if colorAble {
		stdout = colorable.NewColorable(os.Stdout)
	} else {
		stdout = colorable.NewNonColorable(os.Stdout)
	}
	if runtime.GOOS == "windows" {
		io.Copy(stdout, c.StdOut())
	} else {
		io.Copy(stdout, utils.DecoderReader(coder, c.StdOut()))
	}
}

func handleStdErr(c pty.Console) {
	if runtime.GOOS == "windows" {
		io.Copy(os.Stderr, c.StdErr())
	}
}
