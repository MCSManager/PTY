//go:build windows
// +build windows

package winpty

import (
	"fmt"
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
	MouseModes uint
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
}

// the same as open, but uses defaults for Env
func OpenDefault(dllPrefix, cmd, dir string, isColor bool) (*WinPTY, error) {
	var flag uint64 = WINPTY_FLAG_PLAIN_OUTPUT
	if isColor {
		flag = WINPTY_FLAG_COLOR_ESCAPES
	}
	// flag = flag | WINPTY_FLAG_ALLOW_CURPROC_DESKTOP_CREATION
	return OpenWithOptions(Options{
		DllDir:     dllPrefix,
		Command:    cmd,
		Dir:        dir,
		Env:        os.Environ(),
		AgentFlags: flag,
	})
}

func setOptsDefaultValues(options *Options) {
	if options.InitialCols <= 5 {
		options.InitialCols = 50
	}
	if options.InitialRows <= 5 {
		options.InitialRows = 50
	}
	if options.agentTimeoutMs == nil {
		t := uint64(syscall.INFINITE)
		options.agentTimeoutMs = &t
	}
	if options.SpawnFlag != 1 && options.SpawnFlag != 2 {
		options.SpawnFlag = 1
	}
	if options.MouseModes >= 3 {
		options.MouseModes = 0
	}
}

func OpenWithOptions(options Options) (*WinPTY, error) {
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

func (pty *WinPTY) Pid() int {
	pid, _, _ := GetProcessId.Call(pty.procHandle)
	return int(pid)
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
	if pty.Stderr != nil {
		pty.Stderr.Close()
	}
	err := syscall.CloseHandle(syscall.Handle(pty.procHandle))
	if err != nil {
		return err
	}
	pty.closed = true
	return nil

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

func SetMouseMode(winptyConfigT uintptr, mode uint) {
	winpty_config_set_mouse_mode.Call(winptyConfigT, uintptr(mode))
}
