package middlewares

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"

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
				errEmail := services.SendErrorNotification(emailConfig, subject, body)
				if errEmail != nil {
					fmt.Println("Erro ao enviar email:", errEmail)
				}

				errFileName := "error_log.txt"
				ioutil.WriteFile(errFileName, []byte(errMsg), 0644)

				lines := strings.Split(errMsg, "\n")
				firstThreeLines := strings.Join(lines[:3], "\n")

				title := "Erro Interno no Servidor"
				description := "Ocorreu um erro no servidor."
				fields := []models.EmbedField{
					{
						Name:   "Erro Resumido",
						Value:  firstThreeLines,
						Inline: false,
					},
				}

				errWebhook := services.SendDiscordWebhookWithFile(title, description, fields, errFileName)
				if errWebhook != nil {
					fmt.Println("Erro ao enviar webhook:", errWebhook)
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
