package jda

import (
	"encoding/json"
	"reflect"
	"strings"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"fmt"
	"runtime"
	"strconv"
)

func Yay() {
	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		fmt.Printf("%s:%d yay\n", fileName, fileLine)
		return
	}
	fmt.Printf("yay\n")
}

func StartLoggerTimer() {
	for {
		time.Sleep(time.Minute)
		TimestampMutex.Lock()
		Timestamp = time.Now().Format("02/01 15:04")
		TimestampMutex.Unlock()
	}
}

func Getenv(variableName string) string {
	l := GetLogger()
	
	v := os.Getenv(variableName)
	if v == "" {
		l.Error("env variable "+variableName+" not found")
		l.ErrorQueue.DumpErrors()
		os.Exit(1)
	}
	
	return v
}

func GetenvInt(variableName string) int {
	l := GetLogger()
	
	v := Getenv(variableName)
	number, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		l.Error(err.Error())
		l.Error("env variable "+variableName+" is not a valid integer")
		l.ErrorQueue.DumpErrors()
		os.Exit(1)
		return 0
	}
	return int(number)
}

func ProcessCurlyBracketsMacros(
	macroMap map[string]string,
	source string,
) (string, error) {
	l := GetLogger()

	output := ""
	i := 0
	length := len(source)
	for i < length {
		if byte(source[i]) == byte('{') {
			i = i + 1

			macro := ""
			for i < length && byte(source[i]) != byte('}') {
				macro = macro+string(source[i])
				i = i + 1
			}

			if i >= length {
				l.Error("Incomplete macro")
				return "", l.ErrorQueue
			}

			value, ok := macroMap[macro]
			if !ok {
				l.Error("Macro \""+macro+"\" not found")
				return "", l.ErrorQueue
			}
			output = output+value

			i = i + 1
		} else {
			output = output+string(source[i])
			i = i + 1
		}
	}
	return output, nil
}

func Concatenate3Slices(a, b, c []byte) []byte {
	lengthA := len(a)
	lengthB := len(b)
	lengthC := len(c)

	output := make([]byte, lengthA+lengthB+lengthC)

	for i := 0; i < lengthA; i = i + 1 {
		output[i] = a[i]
	}

	tmp := lengthA
	for i := 0; i < lengthB; i = i + 1 {
		output[tmp+i] = b[i]
	}

	tmp = tmp+lengthB
	for i := 0; i < lengthC; i = i + 1 {
		output[tmp+i] = c[i]
	}

	return output
}

func GetStringInBetween(source string, start string, end string) string {
    initIndex := strings.Index(source, start)
    if initIndex == -1 {
        return ""
    }
    initIndex = initIndex + len(start)
    endIndex := strings.Index(source[initIndex:], end)
    if endIndex == -1 {
        return ""
	}
	endIndex = endIndex + initIndex
    return source[initIndex:endIndex]
}

func HttpGetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func GetStructFromFileJsonAndValidate(fileName string, s interface{}) error {
	l := GetLogger()

	fileJson, err := ioutil.ReadFile(fileName)
	if err != nil {
		l.Error(err.Error())
		l.Error("Cannot read file \""+fileName+"\"")
		return l.ErrorQueue
	}

	return GetStructFromJsonAndValidate(fileJson, s)
}

func GetStructFromJsonAndValidate(data []byte, s interface{}) error {
	l := GetLogger()

	err := json.Unmarshal(data, s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Invalid JSON data")
		return l.ErrorQueue
	}

	structValue := reflect.Indirect(reflect.ValueOf(s)) //TODO, make this recursive in validate
	if structValue.Kind() != reflect.Slice {
		isValid, err := Validate(
			structValue.Interface(),
			"in struct validation",
		)
		if !isValid {
			l.Stack(err.(LoggerErrorQueue))
			l.Error("Invalid structure data")
			return l.ErrorQueue
		}
	} else {
		length := structValue.Len()
		fmtIdx := GetFmtIndexer(length)
		for i := 0; i < length; i = i + 1 {
			it := structValue.Index(i).Interface()
			isValid, err := Validate(
				it,
				"in struct validation "+fmtIdx.Format(i),
			)
			if !isValid {
				l.Stack(err.(LoggerErrorQueue))
				l.Error("Invalid structure data")
			}
		}
		if l.ErrorQueue.Queue != nil {
			return l.ErrorQueue
		}
	}

	return nil
}
