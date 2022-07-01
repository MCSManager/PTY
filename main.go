package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/MCSManager/tty/pty"
)

var Dir, Cmd string

func init() {
	flag.StringVar(&Dir, "dir", "", "command work path")
	flag.StringVar(&Cmd, "cmd", "", "command")
}

func main() {
	flag.Parse()
	Pty, err := pty.Start(Dir, Cmd)
	if err != nil {
		fmt.Printf("[MCSMANAGER-TTY] Process Start Error:%s\n", err)
		os.Exit(-1)
	}

	defer func() {
		Pty.Close()
		if err := recover(); err != nil {
			fmt.Printf("[MCSMANAGER-TTY] Recover Point Error: %s", err)
		}
	}()

	Pty.Setsize(50, 50)
	Pty.HandleStdIO()
}
