// Igor Fagundes [ifaresi] 2020
//
// TODO:
//  - only wait if the mutex is unlocked, if don't, await. Making possible 2 or
//    more goroutines see timestamp in same time.

package jda

import (
	"runtime"
	"fmt"
	"strconv"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------
// Logger variables
//-----------------------------------------------------------------------------

var idCounter = uint64(0)
var idCounterMutex sync.Mutex

var Timestamp string = time.Now().Format("02/01 15:04")
var TimestampMutex sync.Mutex

var IsDebug bool = true

var DefaultLoggerDumpCallbackFunc LoggerDumpCallbackFunc = nil

//-----------------------------------------------------------------------------
// Main logger functions
//-----------------------------------------------------------------------------

func getTimestamp() string {
	out := ""
	TimestampMutex.Lock()
	out = Timestamp
	TimestampMutex.Unlock()
	return out
}

type LoggerError struct {
	Content   string
	Stamp     string
	Timestamp string
}

type LoggerErrorQueue struct {
	Queue []LoggerError
}

func (q LoggerErrorQueue) Error() string {
	output := ""
	for _, loggerError := range q.Queue {
		output = output+"\033[0;31m"+loggerError.Timestamp+" "+loggerError.Stamp+
			" ERR\033[0m "+loggerError.Content+"\n"
	}
	q.Queue = nil
	return output
}

func (q *LoggerErrorQueue) Dump(args ...interface{}) {
	for _, loggerError := range q.Queue {
		fmt.Println("\033[0;31m"+loggerError.Timestamp+" "+loggerError.Stamp+
			" ERR\033[0m "+loggerError.Content)
	}
	if DefaultLoggerDumpCallbackFunc != nil {
		DefaultLoggerDumpCallbackFunc(q.Queue, args)
	}
	q.Queue = nil
}

func (q *LoggerErrorQueue) Print() {
	for _, loggerError := range q.Queue {
		fmt.Println("\033[0;31m"+loggerError.Timestamp+" "+loggerError.Stamp+
			" ERR\033[0m "+loggerError.Content)
	}
	q.Queue = nil
}

type Logger struct {
	Location    string
	Stamp       string
	ErrorQueue LoggerErrorQueue
}

type LoggerDumpCallbackFunc func([]LoggerError, ...interface{})

func GetLogger(customLocation ...string) Logger {
	var location string
	if len(customLocation) == 0 {
		//get runtime information
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		location = f.Name()
	} else {
		location = "\""+customLocation[0]+"\""
	}

	//add the track id
	idCounterMutex.Lock()
	stamp := location+" "+strconv.FormatInt(int64(idCounter), 10)
	idCounter = idCounter+1
	idCounterMutex.Unlock()

	return Logger{location, stamp, LoggerErrorQueue{nil}}
}

func (l Logger) NewFrom(locationDetails string) Logger {
	location := l.Location

	//add details, if they exists
	if locationDetails != "" {
		location = location+" "+locationDetails
	}

	//add the track id
	idCounterMutex.Lock()
	stamp := location+" "+strconv.FormatInt(int64(idCounter), 10)
	idCounter = idCounter+1
	idCounterMutex.Unlock()

	return Logger{location, stamp, LoggerErrorQueue{nil}}
}

func (l Logger) Log(content string) {
	fmt.Println("\033[0;32m"+getTimestamp()+" "+l.Stamp+" LOG\033[0m "+content)
}

func (l Logger) Warn(content string) {
	fmt.Println("\033[0;33m"+getTimestamp()+" "+l.Stamp+" WARN\033[0m "+content)
}

func (l *Logger) Error(content string) string {
	l.ErrorQueue.Queue = append(l.ErrorQueue.Queue, LoggerError{content, l.Stamp, getTimestamp()})
	return content
}

func (l Logger) DebugLog(content string) {
	if IsDebug {
		fmt.Println("\033[0;32m"+getTimestamp()+" "+l.Stamp+" LOG DEBUG\033[0m "+content)
	}
}

func (l Logger) DebugWarn(content string) {
	if IsDebug {
		fmt.Println("\033[0;33m"+getTimestamp()+" "+l.Stamp+" WARN DEBUG\033[0m "+content)
	}
}

func (l Logger) DebugError(content string) {
	if IsDebug {
		fmt.Println("\033[0;31m"+getTimestamp()+" "+l.Stamp+" ERR DEBUG\033[0m "+content)
	}
}

func (l *Logger) Stack(errorsQueue LoggerErrorQueue) {
	for _, loggerError := range errorsQueue.Queue {
		l.ErrorQueue.Queue = append(l.ErrorQueue.Queue, loggerError)
	}
}

//-----------------------------------------------------------------------------
// Auxiliar formaters
//-----------------------------------------------------------------------------

//Indexer
type FmtIndexer struct {
	max string
}

func GetFmtIndexer(max int) FmtIndexer {
	return FmtIndexer{strconv.FormatInt(int64(max), 10)}
}

func (f FmtIndexer) Format(i int) string {
	return "["+strconv.FormatInt(int64(i+1), 10)+"/"+f.max+"]"
}

//IP
type FmtIP struct {
	ip string
}

func GetFmtIP(ip string) FmtIP {
	return FmtIP{ip}
}

func FormatIP(ip string) string {
	return "ip"+ip
}

func (f FmtIP) Format() string {
	return "ip"+f.ip
}
