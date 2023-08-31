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

	existingUserByEmail, _ := repository.FindUserByEmail(user.Email)
	if existingUserByEmail != nil {
		return ErrEmailExists
	}

	existingUserByCPF, _ := repository.FindUserByCPF(user.CPF)
	if existingUserByCPF != nil {
		return ErrCPFExists
	}

	hashedCPF, _ := bcrypt.GenerateFromPassword([]byte(formatCPF(user.CPF)), bcrypt.DefaultCost)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.Password = string(hashedPassword)
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
	return len(cpf) == 11
}

func formatCPF(cpf string) string {
	return fmt.Sprintf("%s.%s.%s-%s", cpf[0:3], cpf[3:6], cpf[6:9], cpf[9:11])
}
