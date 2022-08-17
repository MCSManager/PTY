//go:build windows
// +build windows

package winpty

import (
	"fmt"
	"syscall"
	"unsafe"
)

func createAgentCfg(flags uint64) (uintptr, error) {
	var errorPtr uintptr
	defer winpty_error_free.Call(errorPtr)

	winptyConfigT, _, _ := winpty_config_new.Call(uintptr(flags), uintptr(unsafe.Pointer(errorPtr)))
	if winptyConfigT == uintptr(0) {
		return 0, fmt.Errorf("unable to create agent config, %s", GetErrorMessage(errorPtr))
	}

	return winptyConfigT, nil
}

func createSpawnCfg(flags uint32, filePath, cmdline, cwd string, env []string) (uintptr, error) {
	var errorPtr uintptr
	defer winpty_error_free.Call(errorPtr)

	cmdLineStr, err := syscall.UTF16PtrFromString(cmdline)
	if err != nil {
		return 0, fmt.Errorf("failed to convert cmd to pointer")
	}

	filepath, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to convert app name to pointer")
	}

	cwdStr, err := syscall.UTF16PtrFromString(cwd)
	if err != nil {
		return 0, fmt.Errorf("failed to convert working directory to pointer")
	}

	envStr, err := UTF16PtrFromStringArray(env)

	if err != nil {
		return 0, fmt.Errorf("failed to convert cmd to pointer")
	}

	spawnCfg, _, _ := winpty_spawn_config_new.Call(
		uintptr(flags),
		uintptr(unsafe.Pointer(filepath)),
		uintptr(unsafe.Pointer(cmdLineStr)),
		uintptr(unsafe.Pointer(cwdStr)),
		uintptr(unsafe.Pointer(envStr)),
		uintptr(unsafe.Pointer(errorPtr)),
	)

	if spawnCfg == uintptr(0) {
		return 0, fmt.Errorf("unable to create spawn config, %s", GetErrorMessage(errorPtr))
	}

	return spawnCfg, nil
}
