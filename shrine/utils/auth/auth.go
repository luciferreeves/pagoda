package auth

import (
	"shrine/config"
	"shrine/models"
	"shrine/repositories"
	"shrine/utils/crypto"
	"shrine/utils/meta"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	userKey      = "__auth_user"
	tokenHashKey = "__auth_token_hash"
)

func IsAuthenticated(context *fiber.Ctx) bool {
	header, ok := meta.Request(context).Header("Authorization")
	if !ok || !strings.HasPrefix(header, "Bearer ") {
		return false
	}

	rawToken := strings.TrimPrefix(header, "Bearer ")
	tokenHash := crypto.HashToken(rawToken)

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

func RequireStaff(handler fiber.Handler) fiber.Handler {
	return func(context *fiber.Ctx) error {
		if !IsAuthenticated(context) {
			return fiber.ErrUnauthorized
		}
		if !GetUser(context).IsStaff() {
			return fiber.ErrForbidden
		}
		return handler(context)
	}
}

func RequireAdmin(handler fiber.Handler) fiber.Handler {
	return func(context *fiber.Ctx) error {
		if !IsAuthenticated(context) {
			return fiber.ErrUnauthorized
		}
		if !GetUser(context).IsAdmin() {
			return fiber.ErrForbidden
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
	token, err := crypto.GenerateToken()
	if err != nil {
		return "", err
	}

	request := meta.Request(context)
	userAgent, _ := request.Header("User-Agent")

	record := models.Token{
		TokenHash: crypto.HashToken(token),
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