package client

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// HandleShutdown lida com os sinais de desligamento do servidor
func HandleShutdown() {
	// Canal para ouvir sinais de desligamento
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Aguardar o sinal de desligamento
	<-shutdown

	fmt.Println("\nDesligando o servidor...")
	// TODO: Adicionar lógica de limpeza, se necessário
	fmt.Println("Servidor desligado com sucesso.")
}
