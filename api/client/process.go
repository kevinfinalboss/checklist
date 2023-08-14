package client

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
)

func HandleShutdown() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	utils.Logger.Info("\nDesligando o servidor...")
	utils.Logger.Info("Servidor desligado com sucesso.")
}
