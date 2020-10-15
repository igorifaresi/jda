package main

import (
	"context"
	"log"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func RadiusCheckCredentials(secret, name, password, host string) (bool, error) {
	l := GetLogger()

	packet := radius.New(radius.CodeAccessRequest, []byte(secret))
	rfc2865.UserName_SetString(packet, name)
	rfc2865.UserPassword_SetString(packet, password)
	response, err := radius.Exchange(context.TODO(), packet, host)
	if err != nil {
		l.Error(err.Error())
		return false, l.ErrorQueue
	}

	if response.Code != radius.CodeAccessAccept {
		return false, nil
	}

	return true, nil
}
