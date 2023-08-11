package client

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func HandleShutdown() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	fmt.Println("\nDesligando o servidor...")
	fmt.Println("Servidor desligado com sucesso.")
}
