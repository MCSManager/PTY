package start

import (
	"fmt"
	"time"

	pty "github.com/MCSManager/pty/console"

	winio "github.com/Microsoft/go-winio"
)

// \\.\pipe\mypipe
func runControl(fifo string, con pty.Console) error {
	n, err := winio.ListenPipe(fifo, &winio.PipeConfig{})
	if err != nil {
		return fmt.Errorf("open fifo error: %w", err)
	}
	defer n.Close()

	if testFifoResize {
		go func() {
			time.Sleep(time.Second * 5)
			_ = testWinResize(fifo)
		}()
	}

	for {
		conn, err := n.Accept()
		if err != nil {
			return fmt.Errorf("accept fifo error: %w", err)
		}
		go func() {
			defer conn.Close()
			u := newConnUtils(conn, conn)
			_ = handleConn(u, con)
		}()
	}
}

func testWinResize(fifo string) error {
	n, err := winio.DialPipe(fifo, nil)
	if err != nil {
		return fmt.Errorf("open fifo error: %w", err)
	}
	u := newConnUtils(n, n)
	return testResize(u)
}
