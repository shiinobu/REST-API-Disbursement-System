package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"rest-api-disbursement-system/internal/services"
)

const (
	UserIDKey       = "user_id"
	UserUsernameKey = "user_username"
	UserRoleKey     = "user_role"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			errorResponse(c, http.StatusUnauthorized, "Authorization header wajib diisi", gin.H{"authorization": "gunakan format Bearer token"})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			errorResponse(c, http.StatusUnauthorized, "Format token tidak valid", gin.H{"authorization": "gunakan format Bearer token"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &services.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			errorResponse(c, http.StatusUnauthorized, "Token tidak valid", gin.H{"token": "token tidak valid atau sudah kedaluwarsa"})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserUsernameKey, claims.Username)
		c.Set(UserRoleKey, claims.Role)
		c.Next()
	}
}

func errorResponse(c *gin.Context, statusCode int, message string, errors any) {
	if errors == nil {
		errors = gin.H{}
	}

	c.JSON(statusCode, gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}
