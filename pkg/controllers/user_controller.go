package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

func GetUserByCPF(c *gin.Context) {
	cpf := c.Param("cpf")
	user, err := services.GetUserByCPF(cpf)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}
