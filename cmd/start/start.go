package start

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	pty "github.com/MCSManager/pty/console"
	"github.com/MCSManager/pty/utils"
)

var (
	dir, cmd, coder, ptySize, pid, mode  string
	cmds                                 []string
	colorAble, exhaustive, skipExistFile bool
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

	flag.BoolVar(&colorAble, "color", false, "colorable (default false)")
	flag.BoolVar(&skipExistFile, "s", false, "Skip Exist File (default false)")
	flag.BoolVar(&exhaustive, "e", false, "Zip Exhaustive (default false)")
	flag.StringVar(&coder, "coder", "auto", "Coder")
	flag.StringVar(&pid, "pid", "0", "detect pid info")
	flag.StringVar(&dir, "dir", ".", "command work path")
	flag.StringVar(&ptySize, "size", "80,50", "Initialize pty size, stdin will be forwarded directly")
	flag.StringVar(&mode, "m", "pty", "set mode")
}

func Main() {
	flag.Parse()
	args := flag.Args()
	switch mode {
	case "zip":
		if err := utils.Zip(args[:len(args)-1], args[len(args)-1], utils.ZipCfg{Exhaustive: exhaustive}); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "unzip":
		if err := utils.Unzip(args[0], args[1], utils.UnzipCfg{CoderTypes: utils.CoderToType(coder), SkipExistFile: skipExistFile, Exhaustive: exhaustive}); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	case "info":
		runtime.GOMAXPROCS(2)
		info := utils.NewInfo()
		upid, err := strconv.ParseInt(pid, 10, 32)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		utils.Detect(int32(upid), info)
		pinfo, err := json.Marshal(info)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println(string(pinfo))
	default:
		runtime.GOMAXPROCS(6)
		runPTY()
	}
}

func runPTY() {
	json.Unmarshal([]byte(cmd), &cmds)
	con := pty.New(utils.CoderToType(coder))
	if err := con.ResizeWithString(ptySize); err != nil {
		fmt.Printf("[MCSMANAGER-PTY] PTY ReSize Error: %v\n", err)
		return
	}
	err := con.Start(dir, cmds)
	info, _ := json.Marshal(&PtyInfo{
		Pid: con.Pid(),
	})
	fmt.Println(string(info))
	if err != nil {
		fmt.Printf("[MCSMANAGER-PTY] Process Start Error: %v\n", err)
		return
	}
	defer con.Close()
	handleStdIO(con)
	con.Wait()
}

func handleStdIO(c pty.Console) {
	go io.Copy(c.StdIn(), os.Stdin)
	if runtime.GOOS == "windows" && c.StdErr() != nil {
		var stdErr io.Writer
		if colorAble {
			stdErr = pty.Colorable(os.Stderr)
		} else {
			stdErr = pty.NonColorable(os.Stderr)
		}
		go io.Copy(stdErr, c.StdErr())
	}
	handleStdOut(c)
}

func handleStdOut(c pty.Console) {
	var stdout io.Writer
	if colorAble {
		stdout = pty.Colorable(os.Stdout)
	} else {
		stdout = pty.NonColorable(os.Stdout)
	}
	io.Copy(stdout, c.StdOut())
}
