package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/MCSManager/pty/core"
	mytest "github.com/MCSManager/pty/test"
)

var Dir, Cmd string
var test bool

type PtyInfo struct {
	Pid int `json:"pid"`
}

func init() {
	flag.StringVar(&Dir, "dir", ".", "command work path")
	flag.StringVar(&Cmd, "cmd", "", "command")
	flag.StringVar(&core.PtySize, "size", "50,50", "Initialize pty size, stdin will be forwarded directly")
	flag.BoolVar(&core.ColorAble, "color", false, "colorable (default false)")
	flag.StringVar(&core.Coder, "coder", "UTF-8", "Coder")
	flag.BoolVar(&test, "test", false, "Test whether the system environment is pty compatible")
}

func main() {
	flag.Parse()
	if test {
		mytest.Test()
	}
	cmd := []string{}
	json.Unmarshal([]byte(Cmd), &cmd)

	pty, err := core.Start(Dir, cmd)
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error:%v\n", err)
		os.Exit(-1)
	}
	ptyinfo := PtyInfo{
		Pid: pty.Pid(),
	}
	info, _ := json.Marshal(ptyinfo)
	fmt.Printf("%s\n", info)
	defer pty.Close()

	pty.HandleStdIO()
}
