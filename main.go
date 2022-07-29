package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/MCSManager/pty/core"
	t "github.com/MCSManager/pty/test"
)

var Dir, Cmd string
var test bool

func init() {
	flag.StringVar(&Dir, "dir", "", "command work path (default ./)")
	flag.StringVar(&Cmd, "cmd", "", "command")
	flag.StringVar(&core.PtySize, "size", "", "Initialize pty size, stdin will be forwarded directly (default 50,50)")
	flag.BoolVar(&core.Color, "color", false, "colorable (default false)")
	flag.StringVar(&core.Coder, "coder", "UTF-8", "Coder")
	flag.BoolVar(&test, "test", false, "Test whether the system environment is pty compatible")
}

func main() {
	flag.Parse()
	if test {
		t.Test()
	}
	cmd := []string{}
	json.Unmarshal([]byte(Cmd), &cmd)

	pty, err := core.Start(Dir, cmd)
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error:%v\n", err)
		os.Exit(-1)
	}
	fmt.Printf("{pid:%d}\n\n\n\n", pty.Pid())
	defer pty.Close()

	pty.HandleStdIO()
}
