package crypto

import (
	"crypto/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = alphabet[bytes[i]%byte(len(alphabet))]
	}
	return string(bytes)
}

func SystemRef() string {
	return time.Now().Format("020106") + randomString(8)
}

func Ref() string {
	return randomString(12)
}