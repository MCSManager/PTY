package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/MCSManager/pty/core"
)

var Dir, Cmd string

func init() {
	flag.StringVar(&Dir, "dir", "", "command work path (default ./)")
	flag.StringVar(&Cmd, "cmd", "", "command")
	flag.StringVar(&core.PtySize, "size", "", "Initialize pty size, stdin will be forwarded directly (default 50,50)")
	flag.BoolVar(&core.Color, "color", false, "colorable (default false)")
}

func main() {
	flag.Parse()
	fmt.Printf("[MCSMANAGER-TTY] Original command: %s\n", Cmd)

	// 解析命令参数
	cmd := []string{}
	json.Unmarshal([]byte(Cmd), &cmd)

	Pty, err := core.Start(Dir, cmd)
	if err != nil {
		fmt.Printf("[MCSMANAGER-TTY] Process Start Error:%s\n", err)
		os.Exit(-1)
	}
	defer Pty.Close()

	if core.PtySize == "" {
		Pty.Setsize(50, 50)
	} else {
		Pty.ResizeWindow(&core.PtySize)
	}
	Pty.HandleStdIO()
}
