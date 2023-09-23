package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/repository"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if isValidUser(email, password) {
		token, err := services.GenerateToken(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		go func() {
			emailConfig := services.LoadEmailConfig()
			clientIP := getClientIP(c)
			subject := "Registro de Login"
			message := "Você fez login a partir do IP: " + clientIP
			err = services.SendLoginNotification(emailConfig, subject, message)
			if err != nil {
				fmt.Printf("Erro ao enviar e-mail de notificação: %v\n", err)
			}
		}()

		services.SetCookie(c.Writer, token)
		c.Redirect(http.StatusMovedPermanently, "/home?login_success=true")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login?invalid_credentials=true")
}

func validateUserFields(user *models.User) []string {
	fields := map[string]string{
		"Name":            user.Name,
		"Email":           user.Email,
		"Password":        user.Password,
		"CPF":             user.CPF,
		"BirthDate":       user.BirthDate,
		"TelephoneNumber": user.TelephoneNumber,
		"Address":         user.Address,
	}

	missingFields := []string{}
	for field, value := range fields {
		if value == "" {
			missingFields = append(missingFields, field)
		}
	}
	return missingFields
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos no corpo da requisição: " + err.Error()})
		return
	}

	if user.Password != user.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "As senhas não coincidem"})
		return
	}

	missingFields := validateUserFields(&user)
	if len(missingFields) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campos obrigatórios faltando: " + strings.Join(missingFields, ", ")})
		return
	}

	if err := services.CreateUser(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func getClientIP(c *gin.Context) string {
	clientIP := c.ClientIP()
	return clientIP
}
