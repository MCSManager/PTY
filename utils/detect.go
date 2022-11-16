package utils

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

type Info struct {
	Mem          uint64  `json:"mem"`
	Cpu          float64 `json:"cpu"`
	NumConn      int32   `json:"numConn"`
	IOReadSpeed  uint64  `json:"ioReadSpeed"`
	IOWriteSpeed uint64  `json:"ioWriteSpeed"`
	lock         *sync.Mutex
}

func NewInfo() *Info {
	return &Info{lock: &sync.Mutex{}}
}

func Detect(pid int32, info *Info) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return
	}
	if children, err := p.Children(); err == nil {
		for _, v := range children {
			go Detect(v.Pid, info)
		}
	}
	if conn, err := p.Connections(); err == nil {
		atomic.AddInt32(&info.NumConn, int32(len(conn)))
	}
	if io1, err := p.IOCounters(); err == nil {
		time.Sleep(time.Millisecond * 250)
		if io2, err := p.IOCounters(); err == nil {
			atomic.AddUint64(&info.IOReadSpeed, (io2.ReadBytes-io1.ReadBytes)*4)
			atomic.AddUint64(&info.IOWriteSpeed, (io2.WriteBytes-io1.WriteBytes)*4)
		}
	}
	if mem, err := p.MemoryInfo(); err == nil {
		atomic.AddUint64(&info.Mem, mem.RSS)
	}
	if cpu, err := p.CPUPercent(); err == nil {
		info.lock.Lock()
		info.Cpu += cpu
		info.lock.Unlock()
	}
}
