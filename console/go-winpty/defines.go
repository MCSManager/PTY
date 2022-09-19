//go:build windows
// +build windows

package winpty

import (
	"path/filepath"
	"syscall"
)

const (
	WINPTY_SPAWN_FLAG_AUTO_SHUTDOWN       = 1
	WINPTY_SPAWN_FLAG_EXIT_AFTER_SHUTDOWN = 2

	WINPTY_FLAG_CONERR                         = 0x1
	WINPTY_FLAG_PLAIN_OUTPUT                   = 0x2
	WINPTY_FLAG_COLOR_ESCAPES                  = 0x4
	WINPTY_FLAG_ALLOW_CURPROC_DESKTOP_CREATION = 0x8

	WINPTY_MOUSE_MODE_NONE  = 0
	WINPTY_MOUSE_MODE_AUTO  = 1
	WINPTY_MOUSE_MODE_FORCE = 2
)

var (
	modWinPTY *syscall.LazyDLL
	kernel32  *syscall.LazyDLL
	// Error handling...
	winpty_error_code *syscall.LazyProc
	winpty_error_msg  *syscall.LazyProc
	winpty_error_free *syscall.LazyProc

	// Configuration of a new agent.
	winpty_config_new               *syscall.LazyProc
	winpty_config_free              *syscall.LazyProc
	winpty_config_set_initial_size  *syscall.LazyProc
	winpty_config_set_mouse_mode    *syscall.LazyProc
	winpty_config_set_agent_timeout *syscall.LazyProc

	// Start the agent.
	winpty_open          *syscall.LazyProc
	winpty_agent_process *syscall.LazyProc

	// I/O Pipes
	winpty_conin_name  *syscall.LazyProc
	winpty_conout_name *syscall.LazyProc
	winpty_conerr_name *syscall.LazyProc

	// Agent RPC Calls
	winpty_spawn_config_new  *syscall.LazyProc
	winpty_spawn_config_free *syscall.LazyProc
	winpty_spawn             *syscall.LazyProc
	winpty_set_size          *syscall.LazyProc
	winpty_free              *syscall.LazyProc

	//windows api
	GetProcessId *syscall.LazyProc
)

func init() {
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	GetProcessId = kernel32.NewProc("GetProcessId")
}

func setupDefines(dllPrefix string) {

	if modWinPTY != nil {
		return
	}

	modWinPTY = syscall.NewLazyDLL(filepath.Join(dllPrefix, `winpty.dll`))
	// Error handling...
	winpty_error_code = modWinPTY.NewProc("winpty_error_code")
	winpty_error_msg = modWinPTY.NewProc("winpty_error_msg")
	winpty_error_free = modWinPTY.NewProc("winpty_error_free")

	// Configuration of a new agent.
	winpty_config_new = modWinPTY.NewProc("winpty_config_new")
	winpty_config_free = modWinPTY.NewProc("winpty_config_free")
	winpty_config_set_initial_size = modWinPTY.NewProc("winpty_config_set_initial_size")
	winpty_config_set_mouse_mode = modWinPTY.NewProc("winpty_config_set_mouse_mode")
	winpty_config_set_agent_timeout = modWinPTY.NewProc("winpty_config_set_agent_timeout")

	// Start the agent.
	winpty_open = modWinPTY.NewProc("winpty_open")
	winpty_agent_process = modWinPTY.NewProc("winpty_agent_process")

	// I/O Pipes
	winpty_conin_name = modWinPTY.NewProc("winpty_conin_name")
	winpty_conout_name = modWinPTY.NewProc("winpty_conout_name")
	winpty_conerr_name = modWinPTY.NewProc("winpty_conerr_name")

	// Agent RPC Calls
	winpty_spawn_config_new = modWinPTY.NewProc("winpty_spawn_config_new")
	winpty_spawn_config_free = modWinPTY.NewProc("winpty_spawn_config_free")
	winpty_spawn = modWinPTY.NewProc("winpty_spawn")
	winpty_set_size = modWinPTY.NewProc("winpty_set_size")
	winpty_free = modWinPTY.NewProc("winpty_free")
}
