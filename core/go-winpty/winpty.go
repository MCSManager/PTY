//go:build windows
// +build windows

package winpty

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

type Options struct {
	// DllDir is the path to winpty.dll and winpty-agent.exe
	DllDir string
	// FilePath sets the title of the console
	FilePath string
	// Command is the full command to launch
	Command string
	// Dir sets the current working directory for the command
	Dir string
	// Env sets the environment variables. Use the format VAR=VAL.
	Env []string
	// AgentFlags to pass to agent config creation
	AgentFlags uint64
	SpawnFlag  uint32
	MouseModes int
	// Initial size for Columns and Rows
	InitialCols    uint32
	InitialRows    uint32
	agentTimeoutMs *uint64
}

type WinPTY struct {
	Stdin      *os.File
	Stdout     *os.File
	Stderr     *os.File
	pty        uintptr
	procHandle uintptr
	closed     bool
	exitCode   *int
}

func (pty *WinPTY) Pid() int {
	pid, _, _ := GetProcessId.Call(pty.procHandle)
	return int(pid)
}

// 这里不能讲一个file结构体赋值给一个WriteCloser的原因是file结构体的Close方法的第一个参数是一个file指针而不是file指针，也就是说接口方法的对应的结构体方法的第一个参数可以是结构体对或者结构体指针
func (pty *WinPTY) GetStdin() io.Reader {
	return pty.Stdin //这里的类型是机构体指针还是结构体本身取决于该结构体实现的接口方法的第一个参数是指针还是结构体
}

func (pty *WinPTY) GetStdout() io.Writer {
	return pty.Stdout
}

func (pty *WinPTY) GetStderr() io.Writer {
	return pty.Stderr
}

// the same as open, but uses defaults for Env
func OpenPTY(dllPrefix, cmd, dir string, isColor bool) (*WinPTY, error) {
	var flag uint64 = WINPTY_FLAG_PLAIN_OUTPUT
	if isColor {
		flag = WINPTY_FLAG_COLOR_ESCAPES
	}
	flag = flag | WINPTY_FLAG_ALLOW_CURPROC_DESKTOP_CREATION
	return CreateProcessWithOptions(Options{
		DllDir:     dllPrefix,
		Command:    cmd,
		Dir:        dir,
		Env:        os.Environ(),
		AgentFlags: flag,
	})
}

func setOptsDefaultValues(options *Options) {
	// Set the initial size to 40x40 if options is 0
	if options.InitialCols <= 0 {
		options.InitialCols = 40
	}
	if options.InitialRows <= 0 {
		options.InitialRows = 40
	}
	if options.agentTimeoutMs == nil {
		t := uint64(syscall.INFINITE)
		options.agentTimeoutMs = &t
	}
	if options.SpawnFlag == 0 {
		options.SpawnFlag = 1
	}
	if options.MouseModes < 0 {
		options.MouseModes = 0
	}
}

func CreateProcessWithOptions(options Options) (*WinPTY, error) {
	setOptsDefaultValues(&options)
	setupDefines(options.DllDir)
	// create config with specified AgentFlags
	winptyConfigT, err := createAgentCfg(options.AgentFlags)
	if err != nil {
		return nil, err
	}

	winpty_config_set_initial_size.Call(winptyConfigT, uintptr(options.InitialCols), uintptr(options.InitialRows))
	SetMouseMode(winptyConfigT, options.MouseModes)

	var openErr uintptr
	defer winpty_error_free.Call(openErr)
	pty, _, _ := winpty_open.Call(winptyConfigT, uintptr(unsafe.Pointer(openErr)))

	if pty == uintptr(0) {
		return nil, fmt.Errorf("error Launching WinPTY agent, %s", GetErrorMessage(openErr))
	}

	SetAgentTimeout(winptyConfigT, *options.agentTimeoutMs)
	winpty_config_free.Call(winptyConfigT)

	stdinName, _, _ := winpty_conin_name.Call(pty)
	stdoutName, _, _ := winpty_conout_name.Call(pty)
	stderrName, _, _ := winpty_conerr_name.Call(pty)

	obj := &WinPTY{}

	stdinHandle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(stdinName)), syscall.GENERIC_WRITE, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("error getting stdin handle. %s", err)
	}
	obj.Stdin = os.NewFile(uintptr(stdinHandle), "stdin")

	stdoutHandle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(stdoutName)), syscall.GENERIC_READ, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("error getting stdout handle. %s", err)
	}
	obj.Stdout = os.NewFile(uintptr(stdoutHandle), "stdout")

	if options.AgentFlags&WINPTY_FLAG_CONERR == WINPTY_FLAG_CONERR {
		stderrHandle, err := syscall.CreateFile((*uint16)(unsafe.Pointer(stderrName)), syscall.GENERIC_READ, 0, nil, syscall.OPEN_EXISTING, 0, 0)
		if err != nil {
			return nil, fmt.Errorf("error getting stderr handle. %s", err)
		}
		obj.Stderr = os.NewFile(uintptr(stderrHandle), "stderr")
	}

	spawnCfg, err := createSpawnCfg(options.SpawnFlag, options.FilePath, options.Command, options.Dir, options.Env)
	if err != nil {
		return nil, err
	}
	var (
		spawnErr  uintptr
		lastError *uint32
	)
	spawnRet, _, _ := winpty_spawn.Call(pty, spawnCfg, uintptr(unsafe.Pointer(&obj.procHandle)), uintptr(0), uintptr(unsafe.Pointer(lastError)), uintptr(unsafe.Pointer(spawnErr)))
	_, _, _ = winpty_spawn_config_free.Call(spawnCfg)
	defer winpty_error_free.Call(spawnErr)

	if spawnRet == 0 {
		return nil, fmt.Errorf("error spawning process")
	} else {
		obj.pty = pty
		return obj, nil
	}
}

// 设置窗口大小
func (pty *WinPTY) SetSize(wsCol, wsRow uint32) error {
	if wsCol == 0 || wsRow == 0 {
		return fmt.Errorf("wsCol or wsRow = 0")
	}
	_, _, err := winpty_set_size.Call(pty.pty, uintptr(wsCol), uintptr(wsRow), uintptr(0))
	return err
}

// 关闭进程
func (pty *WinPTY) Close() error {
	if pty.closed {
		return nil
	}
	winpty_free.Call(pty.pty)
	pty.Stdin.Close()
	pty.Stdout.Close()
	err := syscall.CloseHandle(syscall.Handle(pty.procHandle))
	if err != nil {
		return err
	}
	pty.closed = true
	return nil

}

func (pty *WinPTY) Wait() error {
	err := WaitForSingleObject(pty.procHandle, syscall.INFINITE)
	pty.Close()
	return err
}

func (pty *WinPTY) GetProcHandle() uintptr {
	return pty.procHandle
}

func (pty *WinPTY) GetAgentProcHandle() uintptr {
	agentProcH, _, _ := winpty_agent_process.Call(pty.pty)
	return agentProcH
}

func SetAgentTimeout(winptyConfigT uintptr, timeoutMs uint64) {
	winpty_config_set_agent_timeout.Call(winptyConfigT, uintptr(timeoutMs))
}

func SetMouseMode(winptyConfigT uintptr, mode int) {
	winpty_config_set_mouse_mode.Call(winptyConfigT, uintptr(mode))
}

func (pty *WinPTY) ExitCode() int {
	if pty.exitCode == nil {
		code, err := GetExitCodeProcess(pty.procHandle)
		if err != nil {
			code = -1
		}
		pty.exitCode = &code
	}
	return *pty.exitCode
}
