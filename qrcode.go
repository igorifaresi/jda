package jda

import qrcode "github.com/skip2/go-qrcode"

func QrcodeGenerate(content string, size int) ([]byte, error) {
	l := GetLogger()
	
	png, err := qrcode.Encode(content, qrcode.Medium, size)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in generate QR Code")
		return []byte{}, l.ErrorQueue
	}
	
	return png, nil
}

func QrcodeGenerateFile(content string, size int, output string) error {
	l := GetLogger()

	err := qrcode.WriteFile(content, qrcode.Medium, size, output)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in generate QR Code and dump to file")
		return l.ErrorQueue
	}

	return nil
}
