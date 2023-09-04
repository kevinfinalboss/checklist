package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWeakPassword     = errors.New("A senha deve ter entre 6 e 20 caracteres, incluir pelo menos uma letra maiúscula, um número e um caractere especial.")
	ErrInvalidCPF       = errors.New("CPF inválido")
	ErrEmailExists      = errors.New("E-mail já cadastrado")
	ErrCPFExists        = errors.New("CPF já cadastrado")
	ErrPasswordMismatch = errors.New("As senhas não coincidem")
)

func CreateUser(user *models.User) error {
	if user.Password != user.ConfirmPassword {
		return ErrPasswordMismatch
	}

	if !isValidPassword(user.Password) {
		return ErrWeakPassword
	}

	if !isValidCPF(user.CPF) {
		return ErrInvalidCPF
	}

	encryptedCPF, err := Encrypt(user.CPF)
	if err != nil {
		return err
	}
	user.CPF = encryptedCPF

	existingUserByEmail, err := repository.FindUserByEmail(user.Email)
	if err == nil && existingUserByEmail != nil {
		return ErrEmailExists
	}

	hash := sha256.Sum256([]byte(user.CPF))
	hashedCPF := hex.EncodeToString(hash[:])

	existingUserByCPF, err := repository.FindUserByCPF(hashedCPF)
	if err == nil && existingUserByCPF != nil {
		return ErrCPFExists
	}

	user.CPF = hashedCPF

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	err = repository.CreateUser(user)
	if err != nil {
		return err
	}

	//smsMessage := fmt.Sprintf("Olá, %s! Sua conta na KevinDev foi criada com sucesso!", user.Name)
	//err = SendSms(user.TelephoneNumber, smsMessage)
	//if err != nil {
	//	return errors.New("Usuário criado, mas falha ao enviar SMS: " + err.Error())
	//}

	currentTime := time.Now()
	formattedTime := currentTime.Format("02/01/2006 às 15:04")

	webhookTitle := "Nova Conta Criada"
	webhookDescription := "Detalhes da nova conta criada:"
	webhookFields := []models.EmbedField{
		{
			Name:   "Nome",
			Value:  user.Name,
			Inline: true,
		},
		{
			Name:   "Email",
			Value:  user.Email,
			Inline: true,
		},
		{
			Name:   "Data de Criação",
			Value:  formattedTime,
			Inline: true,
		},
	}

	err = SendDiscordWebhook(webhookTitle, webhookDescription, webhookFields)
	if err != nil {
		return errors.New("Usuário criado, mas falha ao enviar alerta via webhook: " + err.Error())
	}

	return nil
}

func GetUserByCPF(cpf string) (*models.User, error) {
	users, err := repository.FindAllUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		decryptedCPF, err := Decrypt(user.CPF)
		if err != nil {
			return nil, err
		}

		if decryptedCPF == cpf {
			return &user, nil
		}
	}

	return nil, errors.New("Usuário não encontrado")
}

func GetUserByEmail(email string) (*models.User, error) {
	return repository.FindUserByEmail(email)
}

func GetAllUsers() ([]models.User, error) {
	return repository.FindAllUsers()
}

func isValidPassword(password string) bool {
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasSpecial := strings.ContainsAny(password, "@$!%*?&")
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasUpper && hasLower && hasSpecial && hasNumber
}

func isValidCPF(cpf string) bool {
	return len(cpf) == 11
}
