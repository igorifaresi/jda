package jda

import (
	"reflect"
	"strings"
	"regexp"
	"strconv"
)

var NameRegex = regexp.MustCompile(`^([A-Z]|[a-z])([A-Z]|[a-z]|[-]|[_]|[0-9])*$`)
var CpfRegex = regexp.MustCompile(`^([0-9]{11})$`)
var PhoneRegex = regexp.MustCompile(`^([0-9]{8,15})$`)
var FileNameRegex = regexp.MustCompile(`^([A-Z]|[a-z]|[{]|[.])([A-Z]|[a-z]|[-]|[_]|[:]|[ ]|[0-9]|[.]|[{]|[}])*$`)
var EmailRegex = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)
var SocketRegex = regexp.MustCompile(`^(`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.](`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.](`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.](`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])([:](`+
`[1-9]|[1-9][0-9]|[1-9][0-9][0-9]|[1-9][0-9][0-9][0-9]|`+
`[1-5][0-9][0-9][0-9][0-9]|6[0-4][0-9][0-9][0-9]|`+
`65[0-4][0-9][0-9]|655[0-2][0-9]|6553[0-5]))$`)
var IPRegex = regexp.MustCompile(`^(`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.](`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.](`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])[.](`+
`[0-9]|[1-9][0-9]|1[0-9][0-9]|2[0-4][0-9]|25[0-5])$`)
var FileModeRegex = regexp.MustCompile(`^([0])([0-7])([0-7])([0-7])$`)
var DirectoryRegex = regexp.MustCompile(`(([\w])*([\\]))*$`)

var NumberGTTagRegex = regexp.MustCompile(`^(>)(-)?([0-9])+$`)
var SizeLETagRegex = regexp.MustCompile(`^(size<=)([0-9])+$`)
var MustIfTagRegex = regexp.MustCompile(`^(must-if\()(\w)+(\))([=])(\w|[|])*$`)
var NotMustIfTagRegex = regexp.MustCompile(`^(!must-if\()(\w)+(\))([=])(\w|[|])*$`)

func getStringInBetween(source string, start string, end string) string {
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

type ValidateFieldContext struct {
	V     reflect.Value
	Inter interface{}
	Field reflect.StructField
}

// syntax:
//   must
//
// Makes the field necessary.
func validateMust(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) == "" {
			l.Error("Field \""+ctx.Field.Name+"\" required but was obteined null ("+
				errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field ("+errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   name
//
// Check if the field is a valid name, first letter had to be a latin letter
// (upper or lower case), and each other letter had to be a latin letter
// (upper or lower case) or a number or a '_' or a '-'.
func validateName(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !NameRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a name ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   phone
//
// Check if the field is a valid phone number, without symbols, only phone numbers.
func validatePhone(ctx ValidateFieldContext, errorString string) (bool, error) { //TODO, check CPF sum
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !PhoneRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a phone ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   cpf
//
// Check if the field is a valid CFP, without symbols, only CPF numbers.
func validateCpf(ctx ValidateFieldContext, errorString string) (bool, error) { //TODO, check CPF sum
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !CpfRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a CPF ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   ip
//
// Check if the field is a valid IPv4 adress without port indication.
func validateIP(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !IPRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a ip ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   socket
//
// Check if the field is a valid IPv4 socket.
func validateSocket(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !SocketRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a socket ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   email
//
// Check if the field is a valid e-mail adress.
func validateEmail(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !EmailRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a email ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   file-mode
//
// Check if the field is a valid Unix chmod file permissions octal with 0 in front.
//
// examples:
//   0777
//   0567
//   0123
//   0000
func validateFileMode(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !FileModeRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a file mode ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   directory
//
// Check if the field is a valid unix directory path, the string must to end 
// with a slash ('/') and the paths to the directory had to be expressed with only
// alphanumeric letters.
//
// examples:
//   /path/to/dir/
//   path/to/dir/
//   /
//   dir/
func validateDirectory(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !DirectoryRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a directory ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   file-name
//
// Check if the field is a valid file name, support macros surrounded by curly 
// brackets first letter had to be a latin letter (upper or lower case) or a 
// point or open curly brackets, and each other letter had to be a latin letter
// (upper or lower case) or a number or a '_' or a '-' or a '{' or a '}' or a '.'.
//
// examples:
//   file.zip
//   file123.zip
//   my_executable
//   backup-{actual_date}.conf
func validateFileName(ctx ValidateFieldContext, errorString string) (bool, error) {
	l := GetLogger()

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !FileNameRegex.MatchString(ctx.Inter.(string)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not a valid file name ("+
				errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   >number
//
// Check if the field is a string, and if the size if lower grater than number
// after ">" symbol. Between symbol ">" and the number whitespaces are not allowed.
//
// examples:
//   >1
//   >-1
//   >3
//   >0
func validateNumberGT(ctx ValidateFieldContext, tag, errorString string) (bool, error) {
	l := GetLogger()

	//get the number
	numberString := strings.Split(tag, ">")[1]
	number, _ := strconv.ParseInt(numberString, 10, 64)

	switch ctx.Inter.(type) {
	case int:
		if !(ctx.Inter.(int) > int(number)) {
			l.Error("Field \""+ctx.Field.Name+"\" is not >"+numberString+
				" ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected int ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//   size<=number
//
// Check if the field is a string, and if the size if lower or equals, than number
// after "<=" symbol. The number have to be higher or equals than 1. Between
// keyword "size" and symbol "<=" or between symbol "<=" and the number whitespaces
// are not allowed.
//
// examples:
//   size<=1
//   size<=2
//   size<=300
//   size<=10000
func validateSizeLE(ctx ValidateFieldContext, tag, errorString string) (bool, error) {
	l := GetLogger()

	//get the number
	numberString := strings.Split(tag, "<=")[1]
	number, _ := strconv.ParseInt(numberString, 10, 64)

	switch ctx.Inter.(type) {
	case string:
		if ctx.Inter.(string) != "" && !(int64(len(ctx.Inter.(string))) <= number) {
			l.Error("Field \""+ctx.Field.Name+"\" is not has size<="+numberString+
				" ("+errorString+")")
			return false, l.ErrorQueue
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+"\" field, expected string ("+
			errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//	 must-if(FieldName)=VALUE1|VALUE2
//
// Makes the field necessary if the value of field specified between parentheses
// is equals to any value (separated by vertical bar symbol) specified after
// equals symbol.
func validateMustIf(ctx ValidateFieldContext, tag, errorString string) (bool, error) {
	l := GetLogger()

	//get name of interface to compare
	interName := getStringInBetween(tag, "(", ")")

	//get expression to evaluate, and split the terms of the or expression
	expression := ""
	i := 0
	length := len(tag)
	for i < length {
		if byte(tag[i]) == byte('=') {
			i = i + 1
			for i < length {
				expression = expression+string(tag[i])
				i = i + 1
			}
			break
		}
		i = i + 1
	}
	expressionTerms := strings.Split(expression, "|")

	//get actual interface to compare value
	actualInterValue := ctx.V.FieldByName(interName).Interface()

	//evaluate the expression
	switch actualInterValue.(type) {
	case string:
		for _, term := range expressionTerms {
			if actualInterValue.(string) == term && ctx.Inter.(string) == "" {
				l.Error("Field \""+ctx.Field.Name+"\" required when \""+interName+"\"="+
					"\""+term+"\" ("+errorString+")")
				return false, l.ErrorQueue
			}
		}
	case bool:
		for _, term := range expressionTerms {
			if ((actualInterValue.(bool) && term == "true") ||
			(!actualInterValue.(bool) && term == "false")) && ctx.Inter.(string) == "" {
				l.Error("Field \""+ctx.Field.Name+"\" required when \""+interName+"\"="+
					"\""+term+"\" ("+errorString+")")
				return false, l.ErrorQueue
			}
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+
			"\" field, expected string or bool ("+errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

// syntax:
//	 !must-if(FieldName)=VALUE1|VALUE2
//
// Makes the field not necessary if the value of field specified between
// parentheses is equals to any value (separated by vertical bar symbol)
// specified after equals symbol, but make field necessary if not.
func validateNotMustIf(ctx ValidateFieldContext, tag, errorString string) (bool, error) {
	l := GetLogger()

	//get name of interface to compare
	interName := getStringInBetween(tag, "(", ")")

	//get expression to evaluate, and split the terms of the or expression
	expression := ""
	i := 0
	length := len(tag)
	for i < length {
		if byte(tag[i]) == byte('=') {
			i = i + 1
			for i < length {
				expression = expression+string(tag[i])
				i = i + 1
			}
			break
		}
		i = i + 1
	}
	expressionTerms := strings.Split(expression, "|")

	//get actual interface to compare value
	actualInterValue := ctx.V.FieldByName(interName).Interface()

	//evaluate the expression
	found := false
	value := ""
	switch actualInterValue.(type) {
	case string:
		value = actualInterValue.(string)
		for _, term := range expressionTerms {
			if actualInterValue.(string) == term {
				found = true
				break
			}
		}
	case bool:
		if actualInterValue.(bool) {
			value = "true"
		} else {
			value = "false"
		}
		for _, term := range expressionTerms {
			if ((actualInterValue.(bool) && term == "true") ||
			(!actualInterValue.(bool) && term == "false")) {
				found = true
				break
			}
		}
	default:
		l.Error("Invalid type for \""+ctx.Field.Name+
			"\" field, expected string or bool ("+errorString+")")
		return false, l.ErrorQueue
	}

	//if found the field is not necessary
	if !found && ctx.Inter.(string) == "" {
		l.Error("Field \""+ctx.Field.Name+"\" required when \""+interName+"\"="+
			"\""+value+"\" ("+errorString+")")
		return false, l.ErrorQueue
	}

	return true, nil
}

func Validate(s interface{}, errorString string) (bool, error) {
	l := GetLogger()

	status := true
	v := reflect.ValueOf(s)

	length := v.NumField()
	for i := 0; i < length; i = i + 1 {
		ctx := ValidateFieldContext{}
		ctx.V = v
		ctx.Field = v.Type().Field(i)
		ctx.Inter = v.Field(i).Interface()

		//get validate tags
		body := v.Type().Field(i).Tag.Get("validate")
		if body == "" || body == "-" {
			continue
		}
		tags := strings.Split(body, ",")

		for _, tag := range tags {
			if tag == "must" {
				result, err := validateMust(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "name" {
				result, err := validateName(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "cpf" {
				result, err := validateCpf(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "phone" {
				result, err := validatePhone(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "ip" {
				result, err := validateIP(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "socket" {
				result, err := validateSocket(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "email" {
				result, err := validateEmail(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "file-mode" {
				result, err := validateFileMode(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "directory" {
				result, err := validateDirectory(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else if tag == "file-name" {
				result, err := validateFileName(ctx, errorString)
				if err != nil {
					l.Stack(err.(LoggerErrorQueue))
				}
				if !result {
					status = false
				}
			} else {
				if NumberGTTagRegex.MatchString(tag) {
					result, err := validateNumberGT(ctx, tag, errorString)
					if err != nil {
						l.Stack(err.(LoggerErrorQueue))
					}
					if !result {
						status = false
					}
				} else if SizeLETagRegex.MatchString(tag) {
					result, err := validateSizeLE(ctx, tag, errorString)
					if err != nil {
						l.Stack(err.(LoggerErrorQueue))
					}
					if !result {
						status = false
					}
				} else if MustIfTagRegex.MatchString(tag) {
					result, err := validateMustIf(ctx, tag, errorString)
					if err != nil {
						l.Stack(err.(LoggerErrorQueue))
					}
					if !result {
						status = false
					}
				} else if NotMustIfTagRegex.MatchString(tag) {
					result, err := validateNotMustIf(ctx, tag, errorString)
					if err != nil {
						l.Stack(err.(LoggerErrorQueue))
					}
					if !result {
						status = false
					}
				} else {
					l.DebugError("Unknown validate tag for \""+ctx.Field.Name+
						"\" field ("+errorString+")")
					status = false
				}
			}
		}
	}

	return status, l.ErrorQueue
}