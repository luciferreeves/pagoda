package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"shrine/config"
	"shrine/models"
	"shrine/repositories"
	"shrine/utils/meta"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	userKey      = "__auth_user"
	tokenHashKey = "__auth_token_hash"
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

func IsAuthenticated(context *fiber.Ctx) bool {
	header, ok := meta.Request(context).Header("Authorization")
	if !ok || !strings.HasPrefix(header, "Bearer ") {
		return false
	}

	rawToken := strings.TrimPrefix(header, "Bearer ")
	tokenHash := HashToken(rawToken)

	token, err := repositories.FindValidToken(tokenHash)
	if err != nil {
		return false
	}

	if !token.User.CanAuthenticate() {
		return false
	}

	now := time.Now()
	if token.User.LastSeenAt == nil || time.Since(*token.User.LastSeenAt) > time.Minute {
		token.User.LastSeenAt = &now
		repositories.UpdateLastSeen(&token.User)
	}

	context.Locals(userKey, &token.User)
	context.Locals(tokenHashKey, tokenHash)
	return true
}

func RequireAuthentication(handler fiber.Handler) fiber.Handler {
	return func(context *fiber.Ctx) error {
		if !IsAuthenticated(context) {
			return fiber.ErrUnauthorized
		}
		return handler(context)
	}
}

func GetUser(context *fiber.Ctx) *models.User {
	user, _ := context.Locals(userKey).(*models.User)
	return user
}

func GetTokenHash(context *fiber.Ctx) string {
	hash, _ := context.Locals(tokenHashKey).(string)
	return hash
}

func IssueToken(context *fiber.Ctx, userID uint) (string, error) {
	token, err := GenerateToken()
	if err != nil {
		return "", err
	}

	request := meta.Request(context)
	userAgent, _ := request.Header("User-Agent")

	record := models.Token{
		TokenHash: HashToken(token),
		UserID:    userID,
		ExpiresAt: time.Now().Add(config.Server.TokenExpiry),
		IPAddress: request.IP,
		UserAgent: userAgent,
	}

	if err := repositories.CreateToken(&record); err != nil {
		return "", err
	}

	return token, nil
}