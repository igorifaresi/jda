package waiter

import (
	"net/http"
	"io/ioutil"
	"github.com/igorifaresi/jda"
)

type Context struct {
	W    http.ResponseWriter
	R    *http.Request
	Data []byte
}

type Dish struct {
	Status int
	Text   string
}

type GETFunc func(Context) Dish
type POSTFunc func(Context) Dish

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

func GetQueryParam(paramName string) (string, error) {
	return "", nil	
}

func GetQueryParamInt(paramName string) (int, error) {
	return 0, nil	
}

func POST(path string, handled POSTFunc) {
	l := jda.GetLogger()

	f := func(w http.ResponseWriter, r *http.Request) {
		if Verbose {
			l.Log(`POST request at "`+path+`" ip `+jda.HttpGetIP(r))
		}
		
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case "POST":
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
			dish := handled(Context{ W: w, R: r, Data: body })
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

func Listen() {
	port := jda.Getenv("PORT")
	http.ListenAndServe(":"+port, nil)	
}
