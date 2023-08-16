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

// @Summary Realizar login do usuário
// @Description Autentica o usuário com base no nome de usuário e senha fornecidos e define um cookie de autenticação
// @Tags Autenticação
// @Produce  json
// @Param   username formData string true "Nome de usuário"
// @Param   password formData string true "Senha"
// @Success 301 {string} string "Redireciona para a página inicial com sucesso no login"
// @Failure 303 {string} string "Redireciona para a página de login com credenciais inválidas"
// @Router /login [post]
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
		c.Redirect(http.StatusMovedPermanently, "/home?login_success=true")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login?invalid_credentials=true")
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
		Expires:  time.Now().Add(1 * time.Minute),
		HttpOnly: true,
		Secure:   true,
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
