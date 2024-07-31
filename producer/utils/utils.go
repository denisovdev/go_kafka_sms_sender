package utils

import (
	"math/rand"
	"strings"
)

func GenerateCode(length int) string {
	chars := []rune("123456789")
	b := new(strings.Builder)
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
