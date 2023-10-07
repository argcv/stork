package cntr

import (
	"math/rand"
)

const (
	CharsetUpperCaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CharsetLowerCaseLetters = "abcdefghijklmnopqrstuvwxyz"
	CharsetLetters          = CharsetUpperCaseLetters + CharsetLowerCaseLetters
	CharsetNumbers          = "0123456789"
	DefaultCharset          = CharsetLetters + CharsetNumbers
)

func RandomStringWithCharset(length int, charset string) string {
	if len(charset) == 0 || length < 1 {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return RandomStringWithCharset(length, DefaultCharset)
}
