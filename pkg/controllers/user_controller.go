package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

func GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := services.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetUserAll(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Nenhum usuário encontrado"})
		return
	}
	c.JSON(http.StatusOK, users)
}
