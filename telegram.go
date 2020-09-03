package jda

import (
	"net/http"
	"io/ioutil"
	"strconv"
)

func TelegramSendMessage(botToken string, chatId string, message string) error {
	l := GetLogger()

	resp, err := http.Get("https://api.telegram.org/bot"+botToken+"/sendMessage?"+
		"chat_id="+chatId+
		"&text="+message)
	if err != nil {
		l.Error(err.Error())
		l.Error("Unable to send telegram api http request")
		return l.ErrorQueue
	}

	if resp.StatusCode != 200 {
		l.Error("The request return not OK status code \""+
			strconv.FormatInt(int64(resp.StatusCode), 10)+"\"")
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			l.Error(err.Error())
			l.Error("Cannot read request response body")
			return l.ErrorQueue
		}
		l.Error(string(body))
		return l.ErrorQueue
	}
	return nil
}