package controllers

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/repository"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
	"golang.org/x/crypto/bcrypt"
)

const (
	ErrGeneratingToken    = "Erro ao gerar token"
	ErrInvalidCredentials = "Credenciais inválidas"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if isValidUser(email, password) {
		token, err := generateToken(email)
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

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	missingFields := []string{}
	fieldMap := map[string]string{
		"Name":      user.Name,
		"Password":  user.Password,
		"Email":     user.Email,
		"CPF":       user.CPF,
		"BirthDate": user.BirthDate,
		"Address":   user.Address,
	}

	for fieldName, value := range fieldMap {
		if value == "" {
			missingFields = append(missingFields, fieldName)
		}
	}

	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obrigatórios faltando: " + strings.Join(missingFields, ", ")})
		return
	}

	if err := services.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar usuário"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuário criado com sucesso"})
}

func isValidUser(email, password string) bool {
	user, err := repository.FindUserByEmail(email)
	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
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

func generateToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
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
