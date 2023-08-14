package main

import (
	"os"
	"path/filepath"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kevinfinalboss/checklist-apps/api/client"
	connection "github.com/kevinfinalboss/checklist-apps/pkg/database"
	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
	"github.com/kevinfinalboss/checklist-apps/router"
)

func main() {
	envPath := filepath.Join(".env")

	if err := godotenv.Load(envPath); err != nil {
		color.Red("Erro ao carregar .env")
	}

	if err := utils.LoadConfig(); err != nil {
		color.Red("Erro ao carregar configurações: %v", err)
		return
	}

	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		color.Yellow("PORT não definida no .env, usando 8080 como padrão.")
		port = "8080"
	}

	asciiFigure := figure.NewFigure("CheckList", "", true)
	asciiFigure.Print()

	color.Cyan("Iniciando o servidor na porta %s...", port)

	go func() {
		if err := r.Run(":" + port); err != nil {
			color.Red("Erro ao iniciar o servidor: %v", err)
		}
	}()

	if err := connection.Connect(); err != nil {
		color.Red("Erro ao conectar ao MongoDB: %v", err)
		return
	}
	defer connection.Disconnect()

	color.Green("Servidor rodando na porta %s", port)

	client.HandleShutdown()
}
