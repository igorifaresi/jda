package jda

import qrcode "github.com/skip2/go-qrcode"

func QrcodeGenerateFile(content string, size int, output string) error {
	l := GetLogger()

	err := qrcode.WriteFile(content, qrcode.Medium, size, output)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in generate QR Code")
		return l.ErrorQueue
	}

	return nil
}