package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

// GetUserByEmail godoc
//
//	@Summary		Buscar usuário por e-mail
//	@Description	Busca as informações de um usuário pelo e-mail
//	@Tags			usuário
//	@Accept			json
//	@Produce		json
//	@Router			/user/email/{email} [get]
func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := services.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUserAll godoc
//
//	@Summary		Buscar todos os usuários
//	@Description	Busca todas as informações dos usuários
//	@Tags			usuário
//	@Accept			json
//	@Produce		json
//	@Router			/user/all [get]
func GetUserAll(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum usuário encontrado"})
		return
	}
	c.JSON(http.StatusOK, users)
}
