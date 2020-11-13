package pages

import (
	"math/rand"
	"time"
)

const (
	uniqUnivLen = 16
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var alphaNumericList = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0987654321")
var numericList = []rune("0987654321")

func randRune(runes []rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

func randStringRunes(n int) string {
	return randRune(alphaNumericList, n)
}

// Random returns a random string of length n
func Random(n int) string {
	return randStringRunes(n)
}

// RandomUniqueUniversal returns a universally unique random string (UUID)
func RandomUniqueUniversal() string {
	// TODO: implement UUID
	return Random(uniqUnivLen)
}

// RandomNumeric returns a random number in string format, with n digits
func RandomNumeric(n int) string {
	return randRune(numericList, n)
}
