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
	handleConnection := func(whatcher *SyslogLoggerWatcher, connection net.Conn) {
		output := ""

		for {
			buffer := make([]byte, 128)
			qnt, err := connection.Read(buffer[:])
			goOut := false
			if err != nil {
				goOut = true
			}

			output = output+string(buffer[:qnt])

			if goOut {
				break
			}
		}

		whatcher.ReadMutex.Lock()
		whatcher.Data = whatcher.Data+output
		whatcher.HasNew = true
		whatcher.ReadMutex.Unlock()
	}

	listener, err := net.Listen("tcp", ":4000")
	if err != nil {
			fmt.Println(err)
			return
	}
	defer listener.Close()

	for {
		connection, err := listener.Accept()
        if err != nil {
            continue
		}

		go handleConnection(whatcher, connection)
	}
}