package main

import (
	"fmt"
	"log"
	"unicode"
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

func IsCharLowerCase(char rune) bool {
	return unicode.IsLower(char)
}

func IsCharUpperCase(char rune) bool {
	return unicode.IsUpper(char)
}

func IsCharDigit(char rune) bool {
	return unicode.IsNumber(char)
}

// Password must:
//     - At least one upper case english letter
//     - At least one lower case english letter
//     - At least one digit
//     - Minimum 8 in length
//     - Maximum 20 in length
func IsValidPassword(s string) bool {
	if len(s) < 8 || len(s) > 20 {
		return false
	}

	var (
		runeStr  = []rune(s)
		hasUpper = false
		hasLower = false
		hasDigit = false
	)

	for i := 0; i < len(runeStr); i++ {
		if !hasUpper && IsCharUpperCase(runeStr[i]) {
			hasUpper = true
		} else if !hasLower && IsCharLowerCase(runeStr[i]) {
			hasLower = true
		} else if !hasDigit && IsCharDigit(runeStr[i]) {
			hasDigit = true
		}

		if hasUpper && hasLower && hasDigit {
			return true
		}
	}

	return false
}

func Debug(i ...interface{}) {
	log.Println("[DEBUG] ::", fmt.Sprint(i...))
}
