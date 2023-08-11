package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kevinfinalboss/checklist-apps/api/client"
)

func main() {
	// Caminho para o arquivo .env
	envPath := filepath.Join("..", "..", ".env")

	// Carregar variáveis de ambiente do arquivo .env
	if err := godotenv.Load(envPath); err != nil {
		fmt.Println("Erro ao carregar .env")
	}

	ginMode := os.Getenv("GIN_MODE")
	gin.SetMode(ginMode)

	r := gin.Default()

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
