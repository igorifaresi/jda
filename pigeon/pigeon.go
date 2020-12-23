package pigeon

import (
	"net/http"
	"io/ioutil"
	"bytes"
)

type M map[string]string

type Request struct {
	URL    string
	Query  M
	Header M
	Body   []byte
}

type Response struct {
	Status int
	Text   string
}

func SendRequestPOST(request Request) (Response, error) {
	l := GetLogger()
	
	url := request.URL
	{
		first := true
		for key, element := range request.Query {
			if !first {
				url = url+"&"
			} else {
				url = url+"?"	
			}
			url = url+key+"="+element
		}
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(request.Body))
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
	}
}

func SendRequestGET(request Request) (Response, error) {
	l := GetLogger()
	
	url := request.URL
	{
		first := true
		for key, element := range request.Query {
			if !first {
				url = url+"&"
			} else {
				url = url+"?"	
			}
			url = url+key+"="+element
		}
	}
	
	req, err := http.NewRequest("GET", url, nil)
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
	}
}
