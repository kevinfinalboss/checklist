package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/kevinfinalboss/checklist-apps/pkg/services"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stackTrace := string(debug.Stack())
				errMsg := fmt.Sprintf("ERROR recovered: %v\n%s", err, stackTrace)

				subject := "Erro Interno no Servidor"
				body := "Ocorreu um erro no servidor: " + errMsg
				services.SendEmail(subject, body)

				title := "Erro Interno no Servidor"
				description := "Ocorreu um erro no servidor: " + errMsg
				services.SendDiscordWebhook(title, description)

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
