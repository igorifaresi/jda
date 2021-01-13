package pigeon

import (
	"github.com/igorifaresi/jda"
	"net/http"
	"io/ioutil"
	"bytes"
	"fmt"
)

type M map[string]interface{}

type Request struct {
	Method string
	URL    string
	Query  M
	Header M
	Body   []byte
}

type Response struct {
	Status int
	Text   string
}

type ResponseInter struct {
	Status int
	Inter  interface{}
}

func Fetch(request Request) (Response, error) {
	l := jda.GetLogger()
	
	url := request.URL
	{
		first := true
		for key, element := range request.Query {
			if !first {
				url = url+"&"
			} else {
				url = url+"?"	
			}
			url = url+key+"="+fmt.Sprint(element)
		}
	}
	
	req, err := http.NewRequest(request.Method, url, bytes.NewBuffer(request.Body))
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to create http request")
		return Response{}, l.ErrorQueue
	}
	
	for key, element := range request.Header {
		req.Header.Set(key, element)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to do http request")
		return Response{}, l.ErrorQueue
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to parse response body")
		return Response{}, l.ErrorQueue
	}
	
	return Response{
		Status: resp.StatusCode,
		Text: string(body),
	}, nil
}

func FetchJSON(request Request) (ResponseInter, error) {
	l := jda.GetLogger()
	
	url := request.URL
	{
		first := true
		for key, element := range request.Query {
			if !first {
				url = url+"&"
			} else {
				url = url+"?"	
			}
			url = url+key+"="+fmt.Sprint(element)
		}
	}
	
	req, err := http.NewRequest(request.Method, url, bytes.NewBuffer(request.Body))
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to create http request")
		return ResponseInter{}, l.ErrorQueue
	}
	
	for key, element := range request.Header {
		req.Header.Set(key, element)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to do http request")
		return ResponseInter{}, l.ErrorQueue
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to parse response body")
		return ResponseInter{}, l.ErrorQueue
	}
	
	var inter map[string]interface{}
	err = json.Unmarshal(body, &inter)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to parse response body to JSON")
		return ResponseInter{}, l.ErrorQueue
	}
	
	return ResponseInter{ Status: resp.StatusCode, Inter: inter }, nil
}
