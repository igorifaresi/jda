package jda

import (
	"net/smtp"
)

func SmtpSendMailPlainAuth(
	username string,
	password string,
	host string,
	port string,
	from string,
	to string,
	subject string,
	body string,
) error {
	l := GetLogger()

	auth := smtp.PlainAuth("", username, password, host)

	msg := []byte("To: "+to+"\r\nSubject: "+subject+"\r\n\r\n"+body+"\r\n")

	err := smtp.SendMail(
		host+":"+port,
		auth,
		from,
		[]string{to}, 
		msg,
	)

	if err != nil {
		l.Error(err.Error())
		return l.ErrorQueue
	}

	return nil
}