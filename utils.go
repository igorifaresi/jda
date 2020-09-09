package jda

import (
	"encoding/json"
	"reflect"
)

func CloneInterface(inter interface{}) interface{} {
	value := reflect.ValueOf(inter)
	ty := value.Type()
	var fields []reflect.StructField = nil
	length := value.NumField()
	i := 0
	for i < length {
		fields = append(
			fields,
			reflect.StructField{
				Name: ty.Field(i).Name,
				Type: reflect.TypeOf(ty.Field(i)),
				Tag:  ty.Field(i).Tag,
			},
		)
		i = i + 1
	}
	return reflect.New(reflect.StructOf(fields)).Elem().Interface()
}

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
