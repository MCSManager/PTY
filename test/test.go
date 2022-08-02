package test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/MCSManager/pty/core"
)

var (
	linuxDefault   = []string{"bash", "sh"}
	windowsDefault = []string{"powershell", "cmd"}
)

func Test() {
	var shellPath string
	if runtime.GOOS == "windows" {
		shellPath = lookInPath(windowsDefault)
	} else if runtime.GOOS == "linux" {
		shellPath = lookInPath(linuxDefault)
	}
	if shellPath == "" {
		fmt.Print("0")
		os.Exit(0)
	}
	pty, err := core.Start(".", []string{shellPath})
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error:%v\n", err)
		os.Exit(-1)
	}
	fmt.Print("0")
	pty.Close()
	os.Exit(0)
}

func lookInPath(path []string) string {
	var shellPath string
	for i := 0; i < len(path); i++ {
		shellPath, _ = exec.LookPath(path[i])
		if shellPath != "" {
			continue
		}
	}
	return shellPath
}
