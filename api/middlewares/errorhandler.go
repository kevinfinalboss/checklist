package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("ERROR recovered: %v", err)
				fmt.Println(errMsg)

				subject := "Erro Interno no Servidor"
				body := "Ocorreu um erro no servidor: " + errMsg
				if emailErr := services.SendEmail(subject, body); emailErr != nil {
					fmt.Println("Erro ao enviar e-mail:", emailErr)
				}

				if webhookErr := services.SendDiscordWebhook(errMsg); webhookErr != nil {
					fmt.Println("Erro ao enviar webhook:", webhookErr)
				}

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
