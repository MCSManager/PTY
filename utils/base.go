package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ResizeWindow(sizeText string) (int, int) {
	arr := strings.Split(sizeText, ",")
	if len(arr) != 2 {
		fmt.Printf("[MCSMANAGER-PTY] Set PTY size data failed,original data:%#v\n", sizeText)
		return 50, 50
	}
	cols, err1 := strconv.Atoi(arr[0])
	rows, err2 := strconv.Atoi(arr[1])
	if err1 != nil || err2 != nil {
		fmt.Printf("[MCSMANAGER-PTY] Failed to set window size,original data:%#v\n", sizeText)
		return 50, 50
	}
	return cols, rows
}
