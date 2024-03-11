//go:build !windows
// +build !windows

package start

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	pty "github.com/MCSManager/pty/console"
)

func runControl(fifo string, con pty.Console) error {
	err := os.Remove(fifo)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("remove fifo error: %w", err)
		}
	}
	if err := syscall.Mkfifo(fifo, 0666); err != nil {
		return fmt.Errorf("create fifo error: %w", err)
	}

	if testFifoResize {
		go func() {
			time.Sleep(time.Second * 5)
			_ = testUnixResize(fifo)
		}()
	}

	for {
		f, err := os.OpenFile(fifo, os.O_RDONLY, os.ModeNamedPipe)
		if err != nil {
			return fmt.Errorf("open fifo error: %w", err)
		}
		defer f.Close()
		u := newConnUtils(f, io.Discard)
		_ = handleConn(u, con)
	}
}

func testUnixResize(fifo string) error {
	n, err := os.OpenFile(fifo, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return fmt.Errorf("open fifo error: %w", err)
	}
	defer n.Close()
	u := newConnUtils(nil, n)
	return testResize(u)
}
