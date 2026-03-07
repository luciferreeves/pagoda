package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"shrine/config"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func HashToken(rawToken string) string {
	mac := hmac.New(sha256.New, []byte(config.Server.Secret))
	mac.Write([]byte(rawToken))
	return hex.EncodeToString(mac.Sum(nil))
}