package utils

import "encoding/base64"

func Btoa(i []byte) string {
	return base64.StdEncoding.EncodeToString(i)
}

func Atob(i string) (o []byte) {
	o, err := base64.StdEncoding.DecodeString(i)
	if err != nil {
		return nil
	}
	return
}
