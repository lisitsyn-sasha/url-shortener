package random

import (
	"crypto/rand"
	"math/big"
	"url-shortener/constants"
)

func NewRandomString(size int) string {
	b := make([]rune, size)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(constants.Charset))))
		b[i] = rune(constants.Charset[n.Int64()])
	}
	return string(b)
}
