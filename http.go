package jda

import (
	"net/http"
	"io/ioutil"
	"os"
)

type HttpRequestContext struct {
	W http.ResponseWriter
	R *http.Request
}

type HttpHandleGETFunc func(HttpRequestContext, error) (int, string)
type HttpHandlePOSTFunc func(HttpRequestContext, []byte, error) (int, string)

type HttpMiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

var DefaultHttpMiddlewareFunc HttpMiddlewareFunc = nil

func HttpListenEnvPort() {
	http.ListenAndServe(":"+os.Getenv("PORT"), nil) //Look this better, can have error check
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