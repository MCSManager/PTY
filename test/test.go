package test

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/MCSManager/pty/core"
)

var linuxDefault = []string{"bash", "sh"}

func Test() {
	var cmd []string
	if runtime.GOOS == "windows" {
		cmd = []string{"cmd"}
	} else if runtime.GOOS == "linux" {
		var shellPath string
		for i := 0; i < len(linuxDefault); i++ {
			shellPath, _ = exec.LookPath(linuxDefault[i])
			if shellPath != "" {
				break
			}
		}
		if shellPath == "" {
			fmt.Print("0")
			os.Exit(0)
		}
		cmd = []string{shellPath}
	}
	pty, err := core.Start(".", cmd)
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error:%v\n", err)
		os.Exit(-1)
	}
	fmt.Print("0")
	pty.Close()
	os.Exit(0)
}
