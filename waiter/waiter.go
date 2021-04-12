package waiter

import (
	"net/http"
	"io/ioutil"
	"strconv"
	"github.com/igorifaresi/jda"
)

type Context struct {
	W        http.ResponseWriter
	R        *http.Request
	Data     []byte
	Embedded map[string]interface{}
	Logger   jda.Logger
}

type Dish struct {
	Status int
	Text   string
}

type CallbackFunc func(Context) Dish

const (
	ERROR_NONE = iota
	ERROR_PRINT
	ERROR_DUMP
)
var ErrorMode int = ERROR_PRINT
var Verbose bool = true

func InternalError(text string) Dish {
	return Dish{
		Status: 500,
		Text: "Internal error: "+text,
	}
}

func BadRequest(text string) Dish {
	return Dish{
		Status: 400,
		Text: "Bad request: "+text,
	}
}

func Success(text string) Dish {
	return Dish{
		Status: 200,
		Text: text,
	}
}

func GetQueryParameter(ctx Context, parameterName string) (string, error) {
	l := jda.GetLogger()
	
	values, ok := ctx.R.URL.Query()[parameterName]
	if !ok || len(values) < 1 {
		l.Error(`"`+parameterName+`" query parameter not found`)
		return "", l.ErrorQueue
	}
	return values[0], nil
}

func GetQueryParameterInt(ctx Context, parameterName string) (int, error) {
	l := jda.GetLogger()
	
	values, ok := ctx.R.URL.Query()[parameterName]
	if !ok || len(values) < 1 {
		l.Error(`"`+parameterName+`" query parameter not found`)
		return 0, l.ErrorQueue
	}
	number, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		l.Error(err.Error())
		l.Error(`"`+parameterName+`" query parameter is not a valid integer`)
		return 0, l.ErrorQueue
	}
	return int(number), nil	
}

func GetQueryParameterInt64(ctx Context, parameterName string) (int64, error) {
	l := jda.GetLogger()
	
	values, ok := ctx.R.URL.Query()[parameterName]
	if !ok || len(values) < 1 {
		l.Error(`"`+parameterName+`" query parameter not found`)
		return 0, l.ErrorQueue
	}
	number, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		l.Error(err.Error())
		l.Error(`"`+parameterName+`" query parameter is not a valid 64 bits integer`)
		return 0, l.ErrorQueue
	}
	return number, nil	
}

func GetQueryParameterFloat64(ctx Context, parameterName string) (float64, error) {
	l := jda.GetLogger()
	
	values, ok := ctx.R.URL.Query()[parameterName]
	if !ok || len(values) < 1 {
		l.Error(`"`+parameterName+`" query parameter not found`)
		return 0, l.ErrorQueue
	}
	number, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		l.Error(err.Error())
		l.Error(`"`+parameterName+`" query parameter is not a valid 64 bits float`)
		return 0, l.ErrorQueue
	}
	return number, nil	
}

func GetQueryParameterHex(ctx Context, parameterName string) (int, error) {
	l := jda.GetLogger()
	
	values, ok := ctx.R.URL.Query()[parameterName]
	if !ok || len(values) < 1 {
		l.Error(`"`+parameterName+`" query parameter not found`)
		return 0, l.ErrorQueue
	}
	number, err := strconv.ParseInt(values[0], 16, 64)
	if err != nil {
		l.Error(err.Error())
		l.Error(`"`+parameterName+`" query parameter is not a valid hex integer`)
		return 0, l.ErrorQueue
	}
	return int(number), nil	
}

func do(method string, path string, handled CallbackFunc) {
	f := func(w http.ResponseWriter, r *http.Request) {
		l := jda.GetLogger(path)
		if Verbose {
			l.Log(method+` request at "`+path+`" ip `+jda.HttpGetIP(r))
		}
		
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", method+", OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case method:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				if ErrorMode == ERROR_PRINT || ErrorMode == ERROR_DUMP {
					l.Error("Error in parse request body")
				}
				if ErrorMode == ERROR_PRINT {
					l.ErrorQueue.Print()	
				} else if ErrorMode == ERROR_DUMP {
					l.ErrorQueue.Dump()
				}
				w.WriteHeader(500)
				w.Write([]byte("ifr.waiter: Error in parse request body"))
			}
			dish := handled(Context{ W: w, R: r, Data: body, Logger: l })
			w.WriteHeader(dish.Status)
			w.Write([]byte(dish.Text))
		case "OPTIONS":
			return
		default:
			w.WriteHeader(400)
		}
	}
	http.HandleFunc(path, f)
}

func POST(path string, handled CallbackFunc) {
	do("POST", path, handled)	
}

func GET(path string, handled CallbackFunc) {
	do("GET", path, handled)	
}

type Hostess struct {
	Generator func(CallbackFunc) CallbackFunc 	
}

func NewHostess(f func(CallbackFunc) CallbackFunc) Hostess {
	return Hostess{ Generator: f }
}

func (hostess Hostess) POST(path string, handled CallbackFunc) {
	POST(path, hostess.Generator(handled))
}

func (hostess Hostess) GET(path string, handled CallbackFunc) {
	GET(path, hostess.Generator(handled))
}

func Listen() {
	l := jda.GetLogger()
	
	port := jda.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)
	if Verbose {
		l.Log("Listen at "+port)	
	}
}
