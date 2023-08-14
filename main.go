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
	utils.InitLogger()
	envPath := filepath.Join(".env")

	if err := godotenv.Load(envPath); err != nil {
		utils.LogError(err, "Erro ao carregar .env")
	}

	if err := utils.LoadConfig(); err != nil {
		utils.LogError(err, "Erro ao carregar configurações")
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

	utils.Logger.Infof("Iniciando o servidor na porta %s...", port)

	go func() {
		if err := r.Run(":" + port); err != nil {
			utils.LogError(err, "Erro ao iniciar o servidor")
		}
	}()

	if err := connection.Connect(); err != nil {
		utils.LogError(err, "Erro ao conectar ao MongoDB")
		return
	}
	defer connection.Disconnect()

	utils.Logger.Infof("Servidor rodando na porta %s", port)

	client.HandleShutdown()
}
