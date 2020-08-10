package utils

import "encoding/base64"

func Btoa(i []byte) string {
	return base64.StdEncoding.EncodeToString(i)
}

func Atob(i string) (o []byte, err error) {
	o, err = base64.StdEncoding.DecodeString(i)
	return
}
