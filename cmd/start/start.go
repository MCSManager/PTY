package start

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	pty "github.com/MCSManager/pty/console"
	"github.com/MCSManager/pty/utils"
	"github.com/zijiren233/go-colorable"
	"golang.org/x/term"
)

var (
	dir, cmd, coder, ptySize string
	cmds                     []string
)

type PtyInfo struct {
	Pid int `json:"pid"`
}

func init() {
	if runtime.GOOS == "windows" {
		flag.StringVar(&cmd, "cmd", "[\"cmd\"]", "command")
	} else {
		flag.StringVar(&cmd, "cmd", "[\"sh\"]", "command")
	}

	flag.StringVar(&coder, "coder", "auto", "Coder")
	flag.StringVar(&dir, "dir", ".", "command work path")
	flag.StringVar(&ptySize, "size", "80,50", "Initialize pty size, stdin will be forwarded directly")
}

func Main() {
	flag.Parse()
	runPTY()
}

func runPTY() {
	if err := json.Unmarshal([]byte(cmd), &cmds); err != nil {
		fmt.Println("[MCSMANAGER-PTY] Unmarshal command error: ", err)
		return
	}
	con := pty.New(utils.CoderToType(coder))
	if err := con.ResizeWithString(ptySize); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] PTY Resize error: %v\n", err)
		return
	}
	err := con.Start(dir, cmds)
	info, _ := json.Marshal(&PtyInfo{
		Pid: con.Pid(),
	})
	fmt.Println(string(info))
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process start error: %v\n", err)
		return
	}
	defer con.Close()
	handleStdIO(con)
	_, _ = con.Wait()
}

func handleStdIO(c pty.Console) {
	if colorable.IsReaderTerminal(os.Stdin) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
		go func() { _, _ = io.Copy(c.StdIn(), os.Stdin) }()
	} else {
		go func() { _, _ = io.Copy(c.StdIn(), os.Stdin) }()
	}
	if runtime.GOOS == "windows" && c.StdErr() != nil {
		go func() { _, _ = io.Copy(colorable.NewColorableStderr(), c.StdErr()) }()
	}
	handleStdOut(c)
}

func handleStdOut(c pty.Console) {
	_, _ = io.Copy(colorable.NewColorableStdout(), c.StdOut())
}
