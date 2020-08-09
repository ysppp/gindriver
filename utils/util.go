package utils

import "regexp"

func FilterUsername(str string) bool {
	if ok, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", str); !ok {
		return false
	}
	return true
}
