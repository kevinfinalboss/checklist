package connection

import (
	"context"
	"errors"
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

	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetMaxPoolSize(50)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return errors.New("Erro ao conectar ao MongoDB: " + err.Error())
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return errors.New("Erro ao fazer ping no MongoDB: " + err.Error())
	}

	Client = client
	utils.Logger.Info("Conectado ao MongoDB com sucesso!")
	return nil
}

func Disconnect() {
	if Client != nil {
		Client.Disconnect(context.Background())
		utils.Logger.Info("Desconectado do MongoDB com sucesso!")
	}
}
