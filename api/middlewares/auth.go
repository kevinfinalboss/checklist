package middlewares

import (
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

const (
	ErrUnauthorized = "Unauthorized"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth_token")
		if err != nil || !isValidToken(cookie) {
			c.Redirect(http.StatusSeeOther, "/login?session_expired=true")
			c.Abort()
			return
		}
		c.Next()
	}
}

func isValidToken(cookie string) bool {
	token, err := jwt.ParseWithClaims(cookie, &services.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return false
	}

	if _, ok := token.Claims.(*services.Claims); ok && token.Valid {
		return true
	}

	return false
}
