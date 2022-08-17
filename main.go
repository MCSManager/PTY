package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"

	pty "github.com/MCSManager/pty/core"
)

var dir, cmd, coder, ptySize string
var colorAble, test bool

type PtyInfo struct {
	Pid int `json:"pid"`
}

func init() {
	if runtime.GOOS == "windows" {
		flag.StringVar(&cmd, "cmd", "[\"cmd\"]", "command")
	} else {
		flag.StringVar(&cmd, "cmd", "[\"sh\"]", "command")
	}

	flag.BoolVar(&colorAble, "color", false, "colorable (default false)")
	flag.StringVar(&coder, "coder", "UTF-8", "Coder")
	flag.StringVar(&dir, "dir", ".", "command work path")
	flag.StringVar(&ptySize, "size", "80,50", "Initialize pty size, stdin will be forwarded directly")
	flag.BoolVar(&test, "test", false, "Test whether the system environment is pty compatible")
}

func main() {
	flag.Parse()

	if test {
		fmt.Print("0")
		os.Exit(0)
	}

	con := pty.New(coder, colorAble)

	cmds := []string{}
	json.Unmarshal([]byte(cmd), &cmds)
	if err := con.Start(dir, cmds); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error: %v\n", err)
		os.Exit(1)
	}
	defer con.Close()

	if err := con.ResizeWithString(ptySize); err != nil {
		fmt.Println(err)
	}
	info, _ := json.Marshal(&PtyInfo{
		Pid: con.Pid(),
	})
	fmt.Println(string(info))

	con.HandleStdIO()
	stats, _ := con.Wait()
	os.Exit(stats.ExitCode())
}
