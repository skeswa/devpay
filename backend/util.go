package main

import (
	"fmt"
	"log"
	"regexp"
)

const (
	REGEX_PASSWORD = "^(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$%^&*-]).{8,}$"
)

func String(i interface{}) (string, bool) {
	if i == nil {
		return "", false
	}
	if str, ok := i.(string); ok {
		return str, true
	} else {
		return "", false
	}
}

// Password must:
//     - At least one upper case english letter
//     - At least one lower case english letter
//     - At least one digit
//     - At least one special character
//     - Minimum 8 in length
func IsValidPassword(s string) bool {
	matched, err := regexp.MatchString(REGEX_PASSWORD, s)
	if err != nil {
		Debug("Password was not valid ", err)
		return false
	} else {
		return matched
	}
}

func Debug(i ...interface{}) {
	log.Println("[DEBUG] ::", fmt.Sprint(i...))
}
