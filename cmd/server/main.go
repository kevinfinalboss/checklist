package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kevinfinalboss/checklist-apps/api/client"
	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
	"github.com/kevinfinalboss/checklist-apps/router"
)

func main() {
	envPath := filepath.Join("..", "..", ".env")

	if err := godotenv.Load(envPath); err != nil {
		fmt.Println("Erro ao carregar .env")
	}

	if err := utils.LoadConfig(); err != nil {
		fmt.Println("Erro ao carregar configurações:", err)
		return
	}

	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("PORT não definida no .env, usando 8080 como padrão")
		port = "8080"
	}

	go func() {
		if err := r.Run(":" + port); err != nil {
			fmt.Println("Erro ao iniciar o servidor:", err)
		}
	}()

	fmt.Printf("Servidor rodando na porta %s\n", port)

	client.HandleShutdown()
}
