package connection

import (
	"context"
	"errors"
	"net"
	"os"
	"time"

	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

func Connect() error {
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return errors.New("MONGO_URL não definida no .env")
	}

	dialer := &net.Dialer{}
	dialer.Timeout = 10 * time.Second

	clientOptions := options.Client().
		ApplyURI(mongoURL).
		SetMaxPoolSize(50).
		SetDialer(dialer)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	for i := 0; i < 3; i++ { 
		Client, err = mongo.Connect(ctx, clientOptions)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return errors.New("Erro ao conectar ao MongoDB: " + err.Error())
	}

	err = Client.Ping(ctx, readpref.Primary())
	if err != nil {
		return errors.New("Erro ao fazer ping no MongoDB: " + err.Error())
	}

	utils.Logger.Info("Conectado ao MongoDB com sucesso!")
	return nil
}

func Disconnect() {
	if Client != nil {
		Client.Disconnect(context.Background())
		utils.Logger.Info("Desconectado do MongoDB com sucesso!")
	}
}
