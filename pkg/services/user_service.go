package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWeakPassword = errors.New("A senha deve ter entre 6 e 20 caracteres, incluir pelo menos uma letra maiúscula, um número e um caractere especial.")
	ErrInvalidCPF   = errors.New("CPF inválido")
)

func CreateUser(user *models.User) error {
	if !isValidPassword(user.Password) {
		return ErrWeakPassword
	}

	if !isValidCPF(user.CPF) {
		return ErrInvalidCPF
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	formattedCPF := formatCPF(user.CPF)
	hashedCPF, _ := bcrypt.GenerateFromPassword([]byte(formattedCPF), bcrypt.DefaultCost)
	user.CPF = string(hashedCPF)

	return repository.CreateUser(user)
}

func isValidPassword(password string) bool {
	if len(password) < 6 || len(password) > 20 {
		return false
	}

	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasSpecial := strings.ContainsAny(password, "@$!%*?&")
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasUpper && hasLower && hasSpecial && hasNumber
}

func isValidCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}
	return true
}

func formatCPF(cpf string) string {
	return fmt.Sprintf("%s.%s.%s-%s", cpf[0:3], cpf[3:6], cpf[6:9], cpf[9:11])
}
