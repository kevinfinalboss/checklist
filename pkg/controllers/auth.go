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
			subject := "Registro de Login"
			message := "Você fez login na sua conta. Se não foi você, por favor, entre em contato conosco imediatamente."
			err = services.SendLoginNotification(emailConfig, subject, message)
			if err != nil {
				fmt.Printf("Erro ao enviar e-mail de notificação: %v\n", err)
			}

			title := "Registro de Login"
			description := "Alguém fez login na sua conta. Se não foi você, por favor, entre em contato conosco imediatamente."
			fields := []models.EmbedField{
				{
					Name:   "Email",
					Value:  email,
					Inline: true,
				},
			}
			errWebhook := services.SendDiscordWebhook(title, description, fields)
			if errWebhook != nil {
				fmt.Printf("Erro ao enviar webhook: %v\n", errWebhook)
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
