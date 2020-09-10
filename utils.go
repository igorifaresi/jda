package jda

import (
	"encoding/json"
	"reflect"
)

func GetStructFromJsonAndValidate(data []byte, s interface{}) error {
	l := GetLogger()

	err := json.Unmarshal(data, s)
	if err != nil {
		l.Error(err.Error())
		l.Error("Invalid JSON data")
		return l.ErrorQueue
	}

	isValid, err := Validate(
		reflect.Indirect(reflect.ValueOf(s)).Interface(),
		"in struct validation",
	)
	if !isValid {
		l.Stack(err.(LoggerErrorQueue))
		l.Error("Invalid structure data")
		return l.ErrorQueue
	}

	return nil
}
