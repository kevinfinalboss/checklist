package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type customFormatter struct{}

func (f *customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	date := entry.Time.Format("02/01/2006")
	time := entry.Time.Format("15:04")
	return []byte(fmt.Sprintf("%s: \"%s\" Date: %s ás %s\n", entry.Level.String(), entry.Message, date, time)), nil
}

var Logger *logrus.Logger

func InitLogger() {
	Logger = logrus.New()

	Logger.SetFormatter(&customFormatter{})

	logFileName := fmt.Sprintf("logs/log-%s.log", time.Now().Format("02-01-2006"))
	logFilePath := filepath.Join(logFileName)

	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Logger.Fatalf("Erro ao abrir o arquivo de log: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	Logger.SetOutput(mw)
}

func LogError(err error, message string) {
	if err != nil {
		Logger.Errorf("%s: %v", message, err)
	}
}
