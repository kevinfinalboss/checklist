package controllers

import (
	"github.com/gin-gonic/gin"
)

//	@Summary		Verificar a saúde da API
//	@Description	Retorna OK se a API estiver funcionando
//	@Tags			Health Check
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Router			/diag/health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
