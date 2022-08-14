package console

import (
	"errors"
	"io"
	"os"
	"runtime"

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
	if runtime.GOOS == "windows" {
		io.Copy(c.stdIn(), os.Stdin)
	} else {
		io.Copy(c.stdIn(), utils.EncoderReader(c.coder, os.Stdin))
	}
}

func (c *console) handleStdOut(ColorAble bool) {
	var stdout io.Writer
	if ColorAble {
		stdout = colorable.NewColorable(os.Stdout)
	} else {
		stdout = colorable.NewNonColorable(os.Stdout)
	}
	if runtime.GOOS == "windows" {
		utils.ReadTo("\n", c.stdOut())
		io.Copy(stdout, c.stdOut())
	} else {
		io.Copy(stdout, utils.DecoderReader(c.coder, c.stdOut()))
	}
}
