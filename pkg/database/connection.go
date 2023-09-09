package connection

import (
	"context"
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/kevinfinalboss/checklist-apps/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client  *mongo.Client
	once    sync.Once
	connect sync.Once
)

func Connect() error {
	var err error

	once.Do(func() {
		mongoURL := os.Getenv("MONGO_URL")
		if mongoURL == "" {
			err = errors.New("MONGO_URL não definida no .env")
			return
		}

		dialer := &net.Dialer{}
		dialer.Timeout = 10 * time.Second

		clientOptions := options.Client().
			ApplyURI(mongoURL).
			SetMaxPoolSize(50).
			SetDialer(dialer)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for i := 0; i < 3; i++ {
			Client, err = mongo.Connect(ctx, clientOptions)
			if err == nil {
				break
			}
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			err = errors.New("Erro ao conectar ao MongoDB: " + err.Error())
			return
		}

		err = Client.Ping(ctx, readpref.Primary())
		if err != nil {
			err = errors.New("Erro ao fazer ping no MongoDB: " + err.Error())
			return
		}

		utils.Logger.Info("Conectado ao MongoDB com sucesso!")
	})

	return err
}

func Disconnect() {
	connect.Do(func() {
		if Client != nil {
			Client.Disconnect(context.Background())
			utils.Logger.Info("Desconectado do MongoDB com sucesso!")
		}
	})
}
