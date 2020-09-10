package jda

import (
	"fmt"
	"sync"
	"net"
)

type SyslogLoggerWatcher struct {
	HasNew              bool
	Data                string
	ActualReaderPointer int
	ReadMutex           sync.Mutex
}

func SyslogGetLoggerWhatcher() SyslogLoggerWatcher {
	return SyslogLoggerWatcher{
		HasNew: false,
		Data: "",
		ActualReaderPointer: 0,
	}
}

func SyslogGetLogs(whatcher *SyslogLoggerWatcher) string {
	output := ""
	if whatcher.HasNew {
		whatcher.ReadMutex.Lock()
		output = whatcher.Data[whatcher.ActualReaderPointer:]
		whatcher.ActualReaderPointer = len(whatcher.Data)
		whatcher.HasNew = false
		whatcher.ReadMutex.Unlock()
	}
	return output
}

func SyslogWhatchLogs(whatcher *SyslogLoggerWatcher) {
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
			fmt.Println(err)
			return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
			fmt.Println(err)
			return
	}

	for {
		buffer := make([]byte, 1)
		c.Read(buffer[:])

		whatcher.ReadMutex.Lock()
		whatcher.Data = whatcher.Data+string(buffer)
		whatcher.HasNew = true
		whatcher.ReadMutex.Unlock()
	}
}