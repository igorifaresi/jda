package jda

import (
	"encoding/json"
	"reflect"
	"strings"
	"io/ioutil"
)

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
