package controllers

import (
	"crypto/subtle"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrGeneratingToken    = "Erro ao gerar token"
	ErrInvalidCredentials = "Credenciais inválidas"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if isValidUser(username, password) {
		token, err := generateToken(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrGeneratingToken})
			return
		}

		setCookie(c.Writer, token)
		c.Redirect(http.StatusMovedPermanently, "/home")
		return
	}

	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"error": ErrInvalidCredentials,
	})
}

func isValidUser(username, password string) bool {
	masterUsername := os.Getenv("USER_MASTER")
	masterPassword := os.Getenv("PASSWORD_MASTER")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(masterPassword), bcrypt.DefaultCost)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(username), []byte(masterUsername)) == 1 &&
		bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func generateToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
