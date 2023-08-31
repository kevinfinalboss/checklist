package services

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
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

	existingUserByEmail, err := repository.FindUserByEmail(user.Email)
	if err == nil && existingUserByEmail != nil {
		return ErrEmailExists
	}

	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return err
	}

	hash := sha512.Sum512([]byte(user.CPF + hex.EncodeToString(salt)))
	hashedCPF := hex.EncodeToString(hash[:])

	existingUserByCPF, err := repository.FindUserByCPF(hashedCPF)
	if err == nil && existingUserByCPF != nil {
		return ErrCPFExists
	}

	user.CPF = hashedCPF

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	return repository.CreateUser(user)
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
