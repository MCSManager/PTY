package main

import (
	"io"
	"os"
	"runtime"

	pty "github.com/MCSManager/pty/core"
	"github.com/mattn/go-colorable"
)

func HandleStdIO(c pty.Console) {
	go io.Copy(c.StdIn(), os.Stdin)
	if runtime.GOOS == "windows" {
		go io.Copy(os.Stderr, c.StdErr())
	}
	handleStdOut(c)
}

func handleStdOut(c pty.Console) {
	var stdout io.Writer
	if colorAble {
		stdout = colorable.NewColorable(os.Stdout)
	} else {
		stdout = colorable.NewNonColorable(os.Stdout)
	}
	io.Copy(stdout, c.StdOut())
}
