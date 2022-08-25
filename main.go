package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	pty "github.com/MCSManager/pty/core"
	"github.com/mattn/go-colorable"
)

var dir, cmd, coder, ptySize string
var colorAble bool

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
}

func main() {
	flag.Parse()

	con := pty.New(coder, colorAble)
	if err := con.ResizeWithString(ptySize); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error: %v\n", err)
		return
	}

	cmds := []string{}
	json.Unmarshal([]byte(cmd), &cmds)
	if err := con.Start(dir, cmds); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error: %v\n", err)
		return
	}
	defer con.Close()

	info, _ := json.Marshal(&PtyInfo{
		Pid: con.Pid(),
	})
	fmt.Println(string(info))

	HandleStdIO(con)
	con.Wait()
}

func HandleStdIO(c pty.Console) {
	go io.Copy(c.StdIn(), os.Stdin)
	if runtime.GOOS == "windows" && c.StdErr() != nil {
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
