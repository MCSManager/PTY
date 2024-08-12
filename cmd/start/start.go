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
	fifo                     string
	testFifoResize           bool
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
	flag.StringVar(&fifo, "fifo", "", "control FIFO name")
	flag.BoolVar(&testFifoResize, "test-fifo-resize", false, "test fifo resize")
}

func Main() {
	flag.Parse()
	con, err := newPTY()
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] New pty error: %v\n", err)
		return
	}
	err = con.Start(dir, cmds)
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process start error: %v\n", err)
		return
	}
	info, _ := json.Marshal(&PtyInfo{
		Pid: con.Pid(),
	})
	fmt.Println(string(info))
	defer con.Close()
	if fifo != "" {
		go func() {
			err := runControl(fifo, con)
			if err != nil {
				fmt.Println("[MCSMANAGER-PTY] Control error: ", err)
			}
		}()
	}
	if err = handleStdIO(con); err != nil {
		fmt.Println("[MCSMANAGER-PTY] Handle stdio error: ", err)
	}
	_, _ = con.Wait()
}

func newPTY() (pty.Console, error) {
	if err := json.Unmarshal([]byte(cmd), &cmds); err != nil {
		return nil, fmt.Errorf("unmarshal command error: %w", err)
	}
	con := pty.New(utils.CoderToType(coder))
	if err := con.ResizeWithString(ptySize); err != nil {
		return nil, fmt.Errorf("pty resize error: %w", err)
	}
	return con, nil
}

func handleStdIO(c pty.Console) error {
	if colorable.IsReaderTerminal(os.Stdin) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return fmt.Errorf("make raw error: %w", err)
		}
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
	}
	go func() { _, _ = io.Copy(c.StdIn(), os.Stdin) }()
	if runtime.GOOS == "windows" && c.StdErr() != nil {
		go func() { _, _ = io.Copy(colorable.NewColorableStderr(), c.StdErr()) }()
	}
	_, ok := c.StdOut().(io.WriterTo)
	if !ok {
		return fmt.Errorf("StdOut is not io.WriterTo")
	}
	_, _ = io.Copy(colorable.NewColorableStdout(), c.StdOut())
	return nil
}

const (
	ERROR uint8 = iota + 2
	PING
	RESIZE
)

type errorMsg struct {
	Msg string `json:"msg"`
}

type resizeMsg struct {
	Width  uint `json:"width"`
	Height uint `json:"height"`
}
