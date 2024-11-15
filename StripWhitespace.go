package main

import (
	"strings"
	"unicode"
)

// StripWhitespace removes whitespace from the code example string to more accurately assess whether a code example is a duplicate
// This bypasses issues related to trailing white space or trailing new lines making a code example "appear" unique
func StripWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
