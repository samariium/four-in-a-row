package util

import (
	"crypto/rand"
)

var alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZ")

func NewID(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b)
}
