package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	pty "github.com/MCSManager/pty/core"
	mytest "github.com/MCSManager/pty/test"
)

var dir, cmd, coder, ptySize string
var colorAble, test bool

type PtyInfo struct {
	Pid int `json:"pid"`
}

func init() {
	flag.StringVar(&dir, "dir", ".", "command work path")
	flag.StringVar(&cmd, "cmd", "", "command")
	flag.StringVar(&ptySize, "size", "50,50", "Initialize pty size, stdin will be forwarded directly")
	flag.BoolVar(&colorAble, "color", false, "colorable (default false)")
	flag.StringVar(&coder, "coder", "UTF-8", "Coder")
	flag.BoolVar(&test, "test", false, "Test whether the system environment is pty compatible")
}

func main() {
	flag.Parse()

	if test {
		mytest.Test()
	}

	con := pty.New(coder)

	cmds := []string{}
	json.Unmarshal([]byte(cmd), &cmds)
	if err := con.Start(dir, cmds); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error: %v\n", err)
		os.Exit(-1)
	}
	defer con.Close()
	con.ResizeWithString(ptySize)

	info, _ := json.Marshal(&PtyInfo{
		Pid: con.Pid(),
	})
	fmt.Println(string(info))

	con.HandleStdIO(colorAble)
	con.Wait()
}
