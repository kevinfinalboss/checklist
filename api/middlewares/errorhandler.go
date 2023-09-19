package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := string(debug.Stack())
				errMsg := fmt.Sprintf("ERROR recovered: %v\n%s", err, stackTrace)

				emailConfig := services.LoadEmailConfig()
				subject := "Erro Interno no Servidor"
				body := "Ocorreu um erro no servidor: " + errMsg
				services.SendErrorNotification(emailConfig, subject, body)

				title := "Erro Interno no Servidor"
				description := "Ocorreu um erro no servidor."
				fields := []models.EmbedField{
					{
						Name:   "Erro",
						Value:  fmt.Sprintf("%v", err),
						Inline: false,
					},
					{
						Name:   "Stack Trace",
						Value:  stackTrace,
						Inline: false,
					},
				}
				services.SendDiscordWebhook(title, description, fields)

				c.JSON(http.StatusInternalServerError, gin.H{
					"status": http.StatusInternalServerError,
					"error":  "Internal Server Error",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
