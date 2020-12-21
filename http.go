package jda

import (
	"net/http"
	"io/ioutil"
	"os"
	"strconv"
)

type HttpRequestContext struct {
	W http.ResponseWriter
	R *http.Request
}

type HttpHandleGETFunc func(HttpRequestContext, error) (int, string)
type HttpHandlePOSTFunc func(HttpRequestContext, []byte, error) (int, string)
type HttpHandleWithoutErrorsPOSTFunc func(HttpRequestContext, []byte) (int, string)

type HttpMiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

var DefaultHttpMiddlewareFunc HttpMiddlewareFunc = nil

func HttpListenEnvPort() {
	l := GetLogger()
	
	port := os.Getenv("PORT")
	if port == "" {
		l.Error("env variable PORT is null")
		return
	}
	
	http.ListenAndServe(":"+port, nil) //Look this better, can have error check
}

func HttpHandleGET(path string, handled HttpHandleGETFunc) {
	f := func(w http.ResponseWriter, r *http.Request) {
		ctx := HttpRequestContext{w, r}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case "GET":
			statusCode, output := handled(ctx, nil)

			w.WriteHeader(statusCode)
			w.Write([]byte(output))
		case "OPTIONS":
			return
		default:
			w.WriteHeader(400)
		}
	}
	if DefaultHttpMiddlewareFunc == nil {
		http.HandleFunc(path, f)
	} else {
		http.HandleFunc(path, DefaultHttpMiddlewareFunc(f))
	}
}

func HttpHandlePOST(path string, handled HttpHandlePOSTFunc) {
	f := func(w http.ResponseWriter, r *http.Request) {
		l := GetLogger()

		ctx := HttpRequestContext{w, r}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case "POST":
			var statusCode int
			var output string

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				l.Error("Error in parsing body")
				statusCode, output = handled(ctx, nil, l.ErrorQueue)
			} else {
				statusCode, output = handled(ctx, body, nil)
			}

			w.WriteHeader(statusCode)
			w.Write([]byte(output))
		case "OPTIONS":
			return
		default:
			w.WriteHeader(400)
		}
	}
	if DefaultHttpMiddlewareFunc == nil {
		http.HandleFunc(path, f)
	} else {
		http.HandleFunc(path, DefaultHttpMiddlewareFunc(f))
	}
}

func HttpHandleWithoutErrorsPOST(path string, handled HttpHandleWithoutErrorsPOSTFunc) {
	f := func(w http.ResponseWriter, r *http.Request) {
		l := GetLogger()

		ctx := HttpRequestContext{w, r}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		switch r.Method {
		case "POST":
			var statusCode int
			var output string

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				l.Error("Error in parsing body")
				l.ErrorQueue.PrintErrors()

				statusCode = 500
				output = "Internal error"
			} else {
				statusCode, output = handled(ctx, body)
			}

			w.WriteHeader(statusCode)
			w.Write([]byte(output))
		case "OPTIONS":
			return
		default:
			w.WriteHeader(400)
		}
	}
	if DefaultHttpMiddlewareFunc == nil {
		http.HandleFunc(path, f)
	} else {
		http.HandleFunc(path, DefaultHttpMiddlewareFunc(f))
	}
}

func HttpGetQueryVariable(r *http.Request, variableName string) (string, error) {
	l := GetLogger()
	
	values, ok := r.URL.Query()[variableName]
	if !ok || len(values) < 1 {
		l.Error(`"`+variableName+`" query variable not found`)
		return "", l.ErrorQueue
	}
	return values[0], nil	
}

func HttpGetQueryVariableInt(r *http.Request, variableName string) (int64, error) {
	l := GetLogger()
	
	values, ok := r.URL.Query()[variableName]
	if !ok || len(values) < 1 {
		l.Error(`"`+variableName+`" query variable not found`)
		return 0, l.ErrorQueue
	}
	number, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		l.Error(err.Error())
		l.Error(`"`+variableName+`" query variable is not a valid integer`)
		return 0, l.ErrorQueue
	}
	return number, nil
}
