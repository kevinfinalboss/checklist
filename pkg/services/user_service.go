package services

import (
	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user *models.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	hashedCPF, _ := bcrypt.GenerateFromPassword([]byte(user.CPF), bcrypt.DefaultCost)
	user.CPF = string(hashedCPF)

	return repository.CreateUser(user)
}
