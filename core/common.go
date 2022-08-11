package console

import (
	"errors"
	"io"
	"os"

	"github.com/MCSManager/pty/core/interfaces"
	"github.com/MCSManager/pty/utils"

	"github.com/mattn/go-colorable"
)

var (
	ErrProcessNotStarted = errors.New("process has not been started")
	ErrInvalidCmd        = errors.New("invalid command")
)

type Console interfaces.Console

func New(coder string) Console {
	return newNative(coder)
}

func (c *console) HandleStdIO(ColorAble bool) {
	go c.handleStdIn()
	go c.handleStdOut(ColorAble)
}

func (c *console) handleStdIn() {
	io.Copy(c.stdIn(), utils.Encoder(c.coder, os.Stdin))
}

func (c *console) handleStdOut(ColorAble bool) {
	var stdout io.Writer
	if ColorAble {
		stdout = colorable.NewColorable(os.Stdout)
	} else {
		stdout = colorable.NewNonColorable(os.Stdout)
	}
	io.Copy(stdout, utils.Decoder(c.coder, c.stdOut()))
}
