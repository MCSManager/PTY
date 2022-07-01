package pty

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-colorable"
)

func (pty *Pty) HandleStdIO() {
	go pty.handleStdOut()
	pty.handleStdIn()
}

func (pty *Pty) handleStdIn() {
	inputReader := bufio.NewReader(os.Stdin)
	var getdata Getdata
	var input string
	var cols, rows int
	var err error
	for {
		input, _ = inputReader.ReadString('\n')
		err = json.Unmarshal([]byte(input), &getdata)
		if err != nil {
			fmt.Printf("Unmarshall json err:%v\n,original data:%#v\n", err, input)
			continue
		}
		switch getdata.Type {
		case 1:
			pty.StdIn.Write([]byte(getdata.Data))
		case 2:
			cols, err = strconv.Atoi(strings.Split(getdata.Data, ",")[0])
			if err != nil {
				fmt.Println("SetSize err:", err)
				panic(err)
			}
			rows, err = strconv.Atoi(strings.Split(getdata.Data, ",")[1])
			if err != nil {
				fmt.Println("SetSize err:", err)
				panic(err)
			}
			pty.Setsize(uint32(cols), uint32(rows))
		case 3:
			pty.StdIn.Write([]byte{3})
		default:
			continue
		}
	}
}

func (pty *Pty) handleStdOut() {
	var err error
	var n int
	buf := make([]byte, 2*2048)
	reader := bufio.NewReader(pty.StdOut)
	stdout := colorable.NewColorableStdout()
	for {
		n, err = reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Failed to read from pty master: %s", err)
			continue
		} else if err == io.EOF {
			pty.Close()
			os.Exit(0)
		}
		stdout.Write(buf[:n])
	}
}
