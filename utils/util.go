package utils

import (
	"io/ioutil"
	"regexp"
)

func FilterUsername(str string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", str); !ok {
		return false
	}
	return true
}

func SaveFile(filename string, i []byte) error {
	return ioutil.WriteFile(filename, i, 0644)
}
