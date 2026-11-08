package main

import (
	"strings"
	"unicode"
)

func sanitize(s string, maxLetters int) string {
	s = strings.ToValidUTF8(s, "")

	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, s)

	if r := []rune(s); len(r) > maxLetters {
		s = string(r[:maxLetters])
	}
	return s
}
