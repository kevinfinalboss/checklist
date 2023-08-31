package client

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
)

type ShutdownHandler interface {
	ReleaseResources()
}

func HandleShutdown(handlers ...ShutdownHandler) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	utils.Logger.Info("\nDesligando o servidor...")

	for _, handler := range handlers {
		handler.ReleaseResources()
	}

	utils.Logger.Info("Servidor desligado com sucesso.")
}
