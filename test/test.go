package test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	pty "github.com/MCSManager/pty/core"
)

var (
	linuxDefault   = []string{"bash", "sh"}
	windowsDefault = []string{"powershell", "cmd"}
)

// If the operation is successful, it will output 0 and exit code 0
func Test() {
	var shellPath string
	if runtime.GOOS == "windows" {
		shellPath = lookInPath(windowsDefault)
	} else if runtime.GOOS == "linux" {
		shellPath = lookInPath(linuxDefault)
	}
	console := pty.New("UTF-8")
	console.Start(".", []string{shellPath})
	console.Close()
	fmt.Print("0")
	os.Exit(0)
}

func lookInPath(path []string) string {
	var shellPath string
	for _, v := range path {
		shellPath, _ = exec.LookPath(v)
		if shellPath == "" {
			continue
		}
	}
	return shellPath
}
