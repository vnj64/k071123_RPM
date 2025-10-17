package middleware

import (
	"crypto/rand"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"k071123/internal/shared/permissions"
	"k071123/internal/utils/errs"
	auth "k071123/internal/utils/jwt_helpers"
	"log"
	"math/big"
	"strings"
)

type Middleware struct {
	PublicPemPath string
}

func NewMiddleware(publicPemPath string) *Middleware {
	return &Middleware{PublicPemPath: publicPemPath}
}

// AuthMiddleware проверяет авторизацию с использованием JWT
func (m *Middleware) AuthMiddleware(requiredPermissions []permissions.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf("No Authorization header found")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "auth header is empty",
			))
		}

		tokenParts := strings.Split(authHeader, " ") // Bearer <token>
		if len(tokenParts) != 2 {
			log.Printf("Invalid Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "auth header is invalid",
			))
		}
		tokenString := tokenParts[1]

		claims := jwt.MapClaims{}

		publicKey, err := auth.ParsePublic(m.PublicPemPath)
		if err != nil {
			log.Printf("error %v", err)
			log.Printf("Invalid public key")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "invalid public key",
			))
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		if err != nil || !token.Valid {
			log.Printf("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "invalid token",
			))
		}

		userUUID, ok := claims["sub"].(string)
		if !ok {
			log.Printf("Invalid user UUID")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "invalid user UUID",
			))
		}

		role := permissions.Role(claims["role"].(string))

		c.Locals("uuid", userUUID)
		c.Locals("role", string(role))

		// Проверяем права доступа
		if !hasPermission(role, requiredPermissions) {
			log.Printf("Invalid role permission")
			return c.Status(fiber.StatusForbidden).JSON(errs.NewErrorWithDetails(
				errs.ErrForbidden, "insufficient permissions",
			))
		}

		return c.Next()
	}
}

// AuthMiddlewareWithOptional проверяет авторизацию с использованием JWT, но делает её опциональной
func (m *Middleware) AuthMiddlewareWithOptional(requiredPermissions []permissions.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 {
			log.Printf("Invalid Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "auth header is invalid",
			))
		}
		tokenString := tokenParts[1]

		claims := jwt.MapClaims{}
		publicKey, err := auth.ParsePublic(m.PublicPemPath)
		if err != nil {
			log.Printf("Invalid public key")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "invalid public key",
			))
		}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		if err != nil || !token.Valid {
			log.Printf("Invalid token")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "invalid token",
			))
		}

		userUUID, ok := claims["sub"].(string)
		if !ok {
			log.Printf("Invalid user UUID")
			return c.Status(fiber.StatusUnauthorized).JSON(errs.NewErrorWithDetails(
				errs.ErrUnauthorized, "invalid user UUID",
			))
		}

		role := permissions.Role(claims["role"].(string))
		c.Locals("uuid", userUUID)
		c.Locals("role", string(role))

		if !hasPermission(role, requiredPermissions) {
			log.Printf("Invalid role permission")
			return c.Status(fiber.StatusForbidden).JSON(errs.NewErrorWithDetails(
				errs.ErrForbidden, "insufficient permissions",
			))
		}

		return c.Next()
	}
}

// Функция для проверки наличия разрешений у роли
func hasPermission(role permissions.Role, requiredPermissions []permissions.Permission) bool {
	rolePermissions := permissions.PermissionsForRole(role)

	for _, requiredPermission := range requiredPermissions {
		has := false
		for _, perm := range rolePermissions {
			if perm == requiredPermission {
				has = true
				break
			}
		}
		if !has {
			return false
		}
	}
	return true
}

// GetUserUUIDFromContext извлекает UUID текущего пользователя из контекста (токена)
// Возвращает пустую строку в случае, если пользователь не авторизован.
func GetUserUUIDFromContext(ctx *fiber.Ctx) (string, error) {
	userUUID := ctx.Locals("uuid")
	if userUUID == nil {
		return "", nil
	}
	return userUUID.(string), nil
}

func GetUserRoleFromContext(ctx *fiber.Ctx) (string, error) {
	role := ctx.Locals("role")
	if role == nil {
		return "", errs.NewErrorWithDetails(errs.ErrUnauthorized, "unauthorized: user not found in context")
	}
	return role.(string), nil
}

// GeneratePassword генерирует псевдослучайный пароль
func GeneratePassword() (string, error) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/~"
	passwordLength := 12
	password := make([]byte, passwordLength)

	for i := 0; i < passwordLength; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", err
		}
		password[i] = chars[index.Int64()]
	}

	return string(password), nil
}
