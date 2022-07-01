package main

import (
	"flag"
	"fmt"

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
		fmt.Println("pty start err:", err)
	}
	defer Pty.Close()
	Pty.Setsize(50, 50)
	Pty.HandleStdIO()
}
